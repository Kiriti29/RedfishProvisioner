package redfish


type RedfishBase interface{

    // Inspects hardware according to the Hardware Profile
    InspectHardware(hp string)

    // Sets the hardware profile given by the user. This includes configuring
    // Raids, network devices
    HardwareProfile(hp string)

    // Reads the storage partitioning, Platform config which includes kernel,
    // grub, cpus parameters and creates an iso image by injecting all these
    // parameters into preseed file
    HostProfile(hp string)

    // Triggers the provisioning function which installs OS on baremetal node
    // using the redfish urls and the iso created
    Provision(userdata, hp string) bool

    // Removes the BareMetalHost and all its related resources(secrets,
    // configmaps etc)
    Deprovision(uuid string)

    // Power on a baremetal node
    PowerOn(uuid string)

    // Power off a baremetal node
    PowerOff(uuid string)
}
