package system

import (
	"testing"

	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
)

func TestCalculateCPUPercent(t *testing.T) {
	tests := []struct {
		name     string
		stats    *dockerContainer.StatsResponse
		expected float64
	}{
		{
			name: "zero usage",
			stats: &dockerContainer.StatsResponse{
				Stats: dockerContainer.Stats{
					CPUStats: dockerContainer.CPUStats{
						CPUUsage: dockerContainer.CPUUsage{
							TotalUsage:  1000000000, // 1 second
							PercpuUsage: []uint64{1000000000},
						},
						SystemUsage: 2000000000,
					},
					PreCPUStats: dockerContainer.CPUStats{
						CPUUsage: dockerContainer.CPUUsage{
							TotalUsage: 1000000000,
						},
						SystemUsage: 2000000000,
					},
				},
			},
			expected: 0.0,
		},
		{
			name: "50% usage",
			stats: &dockerContainer.StatsResponse{
				Stats: dockerContainer.Stats{
					CPUStats: dockerContainer.CPUStats{
						CPUUsage: dockerContainer.CPUUsage{
							TotalUsage:  1500000000, // 1.5 seconds
							PercpuUsage: []uint64{1000000000},
						},
						SystemUsage: 3000000000,
					},
					PreCPUStats: dockerContainer.CPUStats{
						CPUUsage: dockerContainer.CPUUsage{
							TotalUsage: 1000000000,
						},
						SystemUsage: 2000000000,
					},
				},
			},
			expected: 50.0,
		},
		{
			name: "100% usage",
			stats: &dockerContainer.StatsResponse{
				Stats: dockerContainer.Stats{
					CPUStats: dockerContainer.CPUStats{
						CPUUsage: dockerContainer.CPUUsage{
							TotalUsage:  2000000000, // 2 seconds
							PercpuUsage: []uint64{1000000000},
						},
						SystemUsage: 3000000000,
					},
					PreCPUStats: dockerContainer.CPUStats{
						CPUUsage: dockerContainer.CPUUsage{
							TotalUsage: 1000000000,
						},
						SystemUsage: 2000000000,
					},
				},
			},
			expected: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateCPUPercent(tt.stats)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateBlkioReadBytes(t *testing.T) {
	tests := []struct {
		name     string
		stats    []dockerContainer.BlkioStatEntry
		expected uint64
	}{
		{
			name: "no read operations",
			stats: []dockerContainer.BlkioStatEntry{
				{Op: "Write", Value: 100},
			},
			expected: 0,
		},
		{
			name: "single read operation",
			stats: []dockerContainer.BlkioStatEntry{
				{Op: "Read", Value: 100},
			},
			expected: 100,
		},
		{
			name: "multiple read operations",
			stats: []dockerContainer.BlkioStatEntry{
				{Op: "Read", Value: 100},
				{Op: "Read", Value: 200},
				{Op: "Write", Value: 50},
			},
			expected: 300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateBlkioReadBytes(tt.stats)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateBlkioWriteBytes(t *testing.T) {
	tests := []struct {
		name     string
		stats    []dockerContainer.BlkioStatEntry
		expected uint64
	}{
		{
			name: "no write operations",
			stats: []dockerContainer.BlkioStatEntry{
				{Op: "Read", Value: 100},
			},
			expected: 0,
		},
		{
			name: "single write operation",
			stats: []dockerContainer.BlkioStatEntry{
				{Op: "Write", Value: 100},
			},
			expected: 100,
		},
		{
			name: "multiple write operations",
			stats: []dockerContainer.BlkioStatEntry{
				{Op: "Write", Value: 100},
				{Op: "Write", Value: 200},
				{Op: "Read", Value: 50},
			},
			expected: 300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateBlkioWriteBytes(tt.stats)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMemoryData(t *testing.T) {
	// This is a basic test to ensure the function returns without error
	// and returns a valid MemoryData struct
	memoryData, err := GetMemoryData()
	assert.NoError(t, err)
	assert.NotNil(t, memoryData)
	assert.Greater(t, memoryData.Total, uint64(0))
	assert.LessOrEqual(t, memoryData.UsedPercent, 100.0)
}

func TestGetCPUData(t *testing.T) {
	// This is a basic test to ensure the function returns without error
	// and returns a valid CPUData struct
	cpuData, err := GetCPUData()
	assert.NoError(t, err)
	assert.NotNil(t, cpuData)
	assert.GreaterOrEqual(t, len(cpuData.Cores), 1)
	assert.LessOrEqual(t, cpuData.Load, 100.0)
}

func TestGetStorageData(t *testing.T) {
	// This is a basic test to ensure the function returns without error
	// and returns a valid StorageData struct
	storageData, err := GetStorageData()
	assert.NoError(t, err)
	assert.NotNil(t, storageData)
	assert.Greater(t, storageData.Total, uint64(0))
	assert.LessOrEqual(t, storageData.UsedPercent, 100.0)
	assert.Greater(t, len(storageData.Disks), 0)
}
