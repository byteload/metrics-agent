package system

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/docker"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
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

	// Calculate actual used memory by subtracting cached, buffers and shared memory
	// This better matches what htop shows as "used" memory
	actualUsed := vm.Used - (vm.Buffers + vm.Cached + vm.Shared)
	actualUsedPercent := float64(actualUsed) / float64(vm.Total) * 100

	return &MemoryData{
		Total:       vm.Total,
		Used:        actualUsed,
		SwapTotal:   swap.Total,
		SwapUsed:    swap.Used,
		UsedPercent: actualUsedPercent,
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
}

// GetDockerContainersData returns information about Docker containers
func GetDockerContainersData() ([]DockerContainerData, error) {
	containers, err := docker.GetDockerStat()
	if err != nil {
		return nil, err
	}

	var containersData []DockerContainerData
	for _, container := range containers {
		containersData = append(containersData, DockerContainerData{
			ID:        container.ContainerID,
			Name:      container.Name,
			Image:     container.Image,
			State:     container.Status,
			Platform:  "linux", // Docker containers typically run on Linux
			StartedAt: "-",     // Not available in current gopsutil version
		})
	}

	return containersData, nil
}

// SystemData represents complete system information
type SystemData struct {
	OS         *OSData               `json:"os"`
	CPU        *CPUData              `json:"cpu"`
	Storage    *StorageData          `json:"storage"`
	Memory     *MemoryData           `json:"memory"`
	Services   []ServiceData         `json:"services"`
	Containers []DockerContainerData `json:"containers"`
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
