import { osInfo, fsSize, mem, services, currentLoad }   from "systeminformation";

// Get OS data
// This function returns an object with the following properties:
// - platform: the platform name
// - distro: the distribution name
// - release: the release version
// - codename: the codename
// - kernel: the kernel version
// - arch: the architecture
export const getOsData = async () => {
    try {
        const osData = await osInfo();
        return {
            platform: osData.platform,
            distro: osData.distro,
            release: osData.release,
            codename: osData.codename,
            kernel: osData.kernel,
            arch: osData.arch,
        };
    } catch (error) {
        console.error("Error getting OS data", error);

        return null;
    }
};

// Get CPU data
// This function returns an object with the following properties:
// - currentLoad: the current load
// - cores: an array with the load of each core
export const getCpuData = async () => {
    try {
        const currentLoadData = await currentLoad();
        return {
            load: currentLoadData.currentLoad,
            cores: currentLoadData.cpus.map(cpu => cpu.load),
        };
    } catch (error) {
        console.error("Error getting CPU data", error);

        return null;
    }
};

// Get storage data
// This function returns an object with the following properties:
// - total: the total storage
// - used: the used storage
// - used_percent: the percentage of used storage
// - disks: an array with the storage data of each disk
export const getStorageData = async () => {
    try {
        const fsData = await fsSize();
        const storageTotal = fsData.reduce((acc, disk) => acc + disk.size, 0);
        const storageUsed = fsData.reduce((acc, disk) => acc + disk.used, 0);

        return {
            total: storageTotal,
            used: storageUsed,
            used_percent: storageUsed / storageTotal * 100,
            disks: fsData.map(disk => {
                return {
                    mount: disk.mount,
                    type: disk.type,
                    size: disk.size,
                    used: disk.used,
                    available: disk.available,
                    use: disk.use,
                };
            }),
        };
    } catch (error) {
        console.error("Error getting storage data", error);

        return null;
    }
};

// Get memory data
// This function returns an object with the following properties:
// - total: the total memory
// - free: the free memory
// - used: the used memory
// - active: the active memory
// - available: the available memory
// - swap_total: the total swap memory
// - swap_used: the used swap memory
// - used_percent: the free swap memory
// - usedPercent: the percentage of used memory
export const getMemoryData = async () => {
    try {
        const memData = await mem();
        return {
            total: memData.total,
            used: memData.total - memData.available,
            swap_total: memData.swaptotal,
            swap_used: memData.swapused,
            used_percent: (memData.total - memData.available) / memData.total * 100,
        };
    } catch (error) {
        console.error("Error getting memory data", error);

        return null;
    }
};

// Get services data
// This function returns an array with the following properties:
// - name: the service name
// - running: the service running status
// - cpu: the service CPU usage
// - mem: the service memory usage
export const getServicesData = async (list) => {
    try {
        return (await services(list)).map(service => {
            return {
                name: service.name,
                running: service.running,
                cpu: service.cpu,
                mem: service.mem,
            };
        });
    } catch (error) {
        console.error("Error getting services data", error);

        return null;
    }
};

// Get system data
// This function returns an object with the following properties:
// - os: the OS data
// - cpu: the CPU data
// - storage: the storage data
// - memory: the memory data
// - services: the services data
export const getSystemData = async ({ services }) => {
    return {
        os: await getOsData(),
        cpu: await getCpuData(),
        storage: await getStorageData(),
        memory: await getMemoryData(),
        services: await getServicesData(services),
    };
};