package system

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	dockerStatsCache     []DockerContainerData
	dockerStatsCacheLock sync.RWMutex
)

// OSData represents operating system information
type OSData struct {
	Platform       string `json:"platform"`
	Distro         string `json:"distro"`
	Release        string `json:"release"`
	Codename       string `json:"codename"`
	Kernel         string `json:"kernel"`
	Arch           string `json:"arch"`
	PlatformFamily string `json:"family"`
}

// GetOSData returns operating system information
func GetOSData() (*OSData, error) {
	info, err := host.Info()
	if err != nil {
		return nil, err
	}

	return &OSData{
		Platform:       info.Platform,
		Distro:         info.PlatformFamily,
		Release:        info.PlatformVersion,
		Kernel:         info.KernelVersion,
		Arch:           info.KernelArch,
		PlatformFamily: info.PlatformFamily,
	}, nil
}

// CPUData represents CPU information
type CPUData struct {
	Load  float64   `json:"load"`
	Cores []float64 `json:"cores"`
}

// GetCPUData returns CPU usage information
func GetCPUData() (*CPUData, error) {
	percent, err := cpu.Percent(0, true)
	if err != nil {
		return nil, err
	}

	totalPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	return &CPUData{
		Load:  totalPercent[0],
		Cores: percent,
	}, nil
}

// DiskInfo represents information about a disk
type DiskInfo struct {
	Mount     string  `json:"mount"`
	Type      string  `json:"type"`
	Size      uint64  `json:"size"`
	Used      uint64  `json:"used"`
	Available uint64  `json:"available"`
	Use       float64 `json:"use"`
}

// StorageData represents storage information
type StorageData struct {
	Total       uint64     `json:"total"`
	Used        uint64     `json:"used"`
	UsedPercent float64    `json:"used_percent"`
	Disks       []DiskInfo `json:"disks"`
}

// GetStorageData returns storage usage information
func GetStorageData() (*StorageData, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var storageData StorageData
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		storageData.Total += usage.Total
		storageData.Used += usage.Used
		storageData.Disks = append(storageData.Disks, DiskInfo{
			Mount:     partition.Mountpoint,
			Type:      partition.Fstype,
			Size:      usage.Total,
			Used:      usage.Used,
			Available: usage.Free,
			Use:       usage.UsedPercent,
		})
	}

	storageData.UsedPercent = float64(storageData.Used) / float64(storageData.Total) * 100
	return &storageData, nil
}

// MemoryData represents memory usage information
type MemoryData struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	SwapTotal   uint64  `json:"swap_total"`
	SwapUsed    uint64  `json:"swap_used"`
	UsedPercent float64 `json:"used_percent"`
}

// GetMemoryData returns memory usage information
func GetMemoryData() (*MemoryData, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	swap, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}

	// Calculate memory usage percentage directly from the virtual memory stats
	// This matches what most system monitors show
	usedPercent := (float64(vm.Used) / float64(vm.Total)) * 100.0

	return &MemoryData{
		Total:       vm.Total,
		Used:        vm.Used,
		SwapTotal:   swap.Total,
		SwapUsed:    swap.Used,
		UsedPercent: usedPercent,
	}, nil
}

// ServiceData represents service information
type ServiceData struct {
	Name    string  `json:"name"`
	Running bool    `json:"running"`
	CPU     float64 `json:"cpu"`
	Mem     float32 `json:"mem"`
}

// GetServicesData returns information about specified services
func GetServicesData(serviceNames []string) ([]ServiceData, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	// Initialize map with all requested services as not running
	serviceMap := make(map[string]*ServiceData)
	for _, name := range serviceNames {
		serviceMap[name] = &ServiceData{
			Name:    name,
			Running: false,
			CPU:     0,
			Mem:     0,
		}
	}

	// Update service data for running processes
	for _, proc := range processes {
		name, err := proc.Name()
		if err != nil {
			continue
		}

		// Check if the process name is in the requested services list
		service, exists := serviceMap[name]
		if !exists {
			continue
		}

		cpu, _ := proc.CPUPercent()
		mem, _ := proc.MemoryPercent()

		// Update service data
		service.Running = true
		service.CPU += cpu
		service.Mem += mem
	}

	// Convert map to slice
	var services []ServiceData
	for _, service := range serviceMap {
		services = append(services, *service)
	}

	return services, nil
}

// DockerContainerData represents Docker container information
type DockerContainerData struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Image     string  `json:"image"`
	State     string  `json:"state"`
	Platform  string  `json:"platform"`
	StartedAt string  `json:"started_at"`
	Restarts  int     `json:"restarts"`
	Ports     []Port  `json:"ports"`
	Mounts    []Mount `json:"mounts"`
	Stats     *Stats  `json:"stats,omitempty"`
}

// Port represents a container port mapping
type Port struct {
	IP          string `json:"ip"`
	PrivatePort uint16 `json:"private_port"`
	PublicPort  uint16 `json:"public_port"`
	Type        string `json:"type"`
}

// Mount represents a container mount point
type Mount struct {
	Type        string `json:"type"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	RW          bool   `json:"rw"`
}

// Stats represents container resource usage statistics
type Stats struct {
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
	MemoryUsage   uint64  `json:"memory_usage"`
	MemoryLimit   uint64  `json:"memory_limit"`
	PIDs          int     `json:"pids"`
	NetworkIO     struct {
		RxBytes   uint64 `json:"rx_bytes"`
		RxPackets uint64 `json:"rx_packets"`
		TxBytes   uint64 `json:"tx_bytes"`
		TxPackets uint64 `json:"tx_packets"`
	} `json:"network_io"`
	BlockIO struct {
		ReadBytes  uint64 `json:"read_bytes"`
		WriteBytes uint64 `json:"write_bytes"`
	} `json:"block_io"`
}

// StartDockerStatsCollector starts collecting Docker stats in the background
func StartDockerStatsCollector(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				stats, err := collectDockerStats()
				if err != nil {
					continue
				}

				dockerStatsCacheLock.Lock()
				dockerStatsCache = stats
				dockerStatsCacheLock.Unlock()
			}
		}
	}()
}

// collectDockerStats collects Docker container stats
func collectDockerStats() ([]DockerContainerData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, dockerContainer.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	// Create a channel to receive container data
	containersChan := make(chan DockerContainerData, len(containers))
	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Process each container in parallel
	for _, container := range containers {
		wg.Add(1)
		go func(container types.Container) {
			defer wg.Done()

			containerData := DockerContainerData{
				ID:        container.ID[:12],
				Name:      container.Names[0][1:],
				Image:     container.Image,
				State:     container.State,
				Platform:  "linux",
				StartedAt: time.Unix(container.Created, 0).Format(time.RFC3339),
				Ports:     make([]Port, 0, len(container.Ports)),
				Mounts:    make([]Mount, 0, len(container.Mounts)),
			}

			for _, p := range container.Ports {
				containerData.Ports = append(containerData.Ports, Port{
					IP:          p.IP,
					PrivatePort: p.PrivatePort,
					PublicPort:  p.PublicPort,
					Type:        p.Type,
				})
			}

			for _, m := range container.Mounts {
				containerData.Mounts = append(containerData.Mounts, Mount{
					Type:        string(m.Type),
					Source:      m.Source,
					Destination: m.Destination,
					Mode:        m.Mode,
					RW:          m.RW,
				})
			}

			// Only fetch stats for running containers
			if container.State == "running" {
				// Create a new client for each goroutine to ensure thread safety
				cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
				if err == nil {
					defer cli.Close()

					stats, err := cli.ContainerStats(ctx, container.ID, false)
					if err == nil {
						var statsJSON dockerContainer.StatsResponse
						if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err == nil {
							memPercent := float64(0)
							if statsJSON.MemoryStats.Limit != 0 {
								memPercent = float64(statsJSON.MemoryStats.Usage) / float64(statsJSON.MemoryStats.Limit) * 100
							}

							containerData.Stats = &Stats{
								CPUPercent:    calculateCPUPercent(&statsJSON),
								MemoryPercent: memPercent,
								MemoryUsage:   statsJSON.MemoryStats.Usage,
								MemoryLimit:   statsJSON.MemoryStats.Limit,
								PIDs:          int(statsJSON.PidsStats.Current),
								NetworkIO: struct {
									RxBytes   uint64 `json:"rx_bytes"`
									RxPackets uint64 `json:"rx_packets"`
									TxBytes   uint64 `json:"tx_bytes"`
									TxPackets uint64 `json:"tx_packets"`
								}{
									RxBytes:   statsJSON.Networks["eth0"].RxBytes,
									RxPackets: statsJSON.Networks["eth0"].RxPackets,
									TxBytes:   statsJSON.Networks["eth0"].TxBytes,
									TxPackets: statsJSON.Networks["eth0"].TxPackets,
								},
								BlockIO: struct {
									ReadBytes  uint64 `json:"read_bytes"`
									WriteBytes uint64 `json:"write_bytes"`
								}{
									ReadBytes:  calculateBlkioReadBytes(statsJSON.BlkioStats.IoServiceBytesRecursive),
									WriteBytes: calculateBlkioWriteBytes(statsJSON.BlkioStats.IoServiceBytesRecursive),
								},
							}
						}
						stats.Body.Close()
					}
				}
			} else {
				// Set zero stats for non-running containers
				containerData.Stats = &Stats{
					CPUPercent:    0,
					MemoryPercent: 0,
					MemoryUsage:   0,
					MemoryLimit:   0,
					PIDs:          0,
					NetworkIO: struct {
						RxBytes   uint64 `json:"rx_bytes"`
						RxPackets uint64 `json:"rx_packets"`
						TxBytes   uint64 `json:"tx_bytes"`
						TxPackets uint64 `json:"tx_packets"`
					}{},
					BlockIO: struct {
						ReadBytes  uint64 `json:"read_bytes"`
						WriteBytes uint64 `json:"write_bytes"`
					}{},
				}
			}

			containersChan <- containerData
		}(container)
	}

	// Start a goroutine to close the channel when all workers are done
	go func() {
		wg.Wait()
		close(containersChan)
	}()

	// Collect results from the channel
	containersData := make([]DockerContainerData, 0, len(containers))
	for containerData := range containersChan {
		containersData = append(containersData, containerData)
	}

	return containersData, nil
}

// GetDockerContainersData returns cached Docker container information
func GetDockerContainersData() ([]DockerContainerData, error) {
	dockerStatsCacheLock.RLock()
	defer dockerStatsCacheLock.RUnlock()

	if dockerStatsCache == nil {
		// If cache is not initialized yet, collect stats immediately
		stats, err := collectDockerStats()
		if err != nil {
			return nil, err
		}
		return stats, nil
	}

	// Return a copy of the cache to prevent data races
	result := make([]DockerContainerData, len(dockerStatsCache))
	copy(result, dockerStatsCache)
	return result, nil
}

// calculateCPUPercent calculates CPU usage percentage
func calculateCPUPercent(stats *dockerContainer.StatsResponse) float64 {
	var cpuPercent float64
	if stats.CPUStats.SystemUsage != 0 {
		cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)

		if systemDelta > 0 && cpuDelta > 0 {
			cpuPercent = (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100
		}
	}
	return cpuPercent
}

// calculateBlkioReadBytes calculates total read bytes from BlkioStats
func calculateBlkioReadBytes(stats []dockerContainer.BlkioStatEntry) uint64 {
	var total uint64
	for _, stat := range stats {
		if stat.Op == "Read" {
			total += stat.Value
		}
	}
	return total
}

// calculateBlkioWriteBytes calculates total write bytes from BlkioStats
func calculateBlkioWriteBytes(stats []dockerContainer.BlkioStatEntry) uint64 {
	var total uint64
	for _, stat := range stats {
		if stat.Op == "Write" {
			total += stat.Value
		}
	}
	return total
}

// SystemData represents complete system information
type SystemData struct {
	OS       *OSData       `json:"os"`
	CPU      *CPUData      `json:"cpu"`
	Storage  *StorageData  `json:"storage"`
	Memory   *MemoryData   `json:"memory"`
	Services []ServiceData `json:"services"`
}

// GetSystemData returns complete system information
func GetSystemData(services []string) (*SystemData, error) {
	osData, err := GetOSData()
	if err != nil {
		return nil, err
	}

	cpuData, err := GetCPUData()
	if err != nil {
		return nil, err
	}

	storageData, err := GetStorageData()
	if err != nil {
		return nil, err
	}

	memoryData, err := GetMemoryData()
	if err != nil {
		return nil, err
	}

	servicesData, err := GetServicesData(services)
	if err != nil {
		return nil, err
	}

	return &SystemData{
		OS:       osData,
		CPU:      cpuData,
		Storage:  storageData,
		Memory:   memoryData,
		Services: servicesData,
	}, nil
}
