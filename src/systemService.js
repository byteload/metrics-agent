import { osInfo, fsSize, mem, services, currentLoad, processes, processLoad }   from "systeminformation";

// Get OS data
// This function returns an object with the following properties:
// - platform: the platform name
// - distro: the distribution name
// - release: the release version
// - codename: the codename
// - kernel: the kernel version
// - arch: the architecture
const getOsData = async () => {
    const osData = await osInfo();
    return {
        platform: osData.platform,
        distro: osData.distro,
        release: osData.release,
        codename: osData.codename,
        kernel: osData.kernel,
        arch: osData.arch,
    };
};

// Get CPU data
// This function returns an object with the following properties:
// - currentLoad: the current load
// - avgLoad: the average load
// - user: the user load
// - system: the system load
// - cores: an array with the load of each core
const getCpuData = async () => {
    const currentLoadData = await currentLoad();
    return {
        currentLoad: currentLoadData.currentLoad,
        avgLoad: currentLoadData.avgLoad,
        user: currentLoadData.currentLoadUser,
        system: currentLoadData.currentLoadSystem,
        cores: currentLoadData.cpus.map(cpu => cpu.load),
    };
};

// Get storage data
// This function returns an object with the following properties:
// - total: the total storage
// - used: the used storage
// - usedPercent: the percentage of used storage
// - disks: an array with the storage data of each disk
const getStorageData = async () => {
    const fsData = await fsSize();
    const storageTotal = fsData.reduce((acc, disk) => acc + disk.size, 0);
    const storageUsed = fsData.reduce((acc, disk) => acc + disk.used, 0);

    return {
        total: storageTotal,
        used: storageUsed,
        usedPercent: storageUsed / storageTotal * 100,
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
};

// Get memory data
// This function returns an object with the following properties:
// - total: the total memory
// - free: the free memory
// - used: the used memory
// - active: the active memory
// - available: the available memory
// - buffcache: the buffer cache memory
// - swaptotal: the total swap memory
// - swapused: the used swap memory
// - swapfree: the free swap memory
// - usedPercent: the percentage of used memory
const getMemoryData = async () => {
    const memData = await mem();
    return {
        total: memData.total,
        free: memData.free,
        used: memData.used,
        active: memData.active,
        available: memData.available,
        buffcache: memData.buffcache,
        swaptotal: memData.swaptotal,
        swapused: memData.swapused,
        swapfree: memData.swapfree,
        usedPercent: memData.active / memData.total * 100,
    };
};

// Get services data
// This function returns an array with the following properties:
// - name: the service name
// - running: the service running status
// - cpu: the service CPU usage
// - mem: the service memory usage
const getServicesData = async (list) => {
    return (await services(list)).map(service => {
        return {
            name: service.name,
            running: service.running,
            cpu: service.cpu,
            mem: service.mem,
        };
    });
};

// Get system data
// This function returns an object with the following properties:
// - os: the OS data
// - cpu: the CPU data
// - storage: the storage data
// - memory: the memory data
// - services: the services data
export const getSystemData = async () => {
    return {
        os: await getOsData(),
        cpu: await getCpuData(),
        storage: await getStorageData(),
        memory: await getMemoryData(),
        services: await getServicesData('mysql,nginx,redis,php,node'),
    };
};