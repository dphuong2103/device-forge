# Privileged Helper Decision

The MVP does not install or run a privileged helper.

The per-user Agent owns its data, Android SDK directory, AVD directories, emulator processes, loopback ports, tray integration, and updates staged inside the user profile. WHPX enablement is a documented host prerequisite performed by the user or an administrator outside the Agent.

If a later installer or machine-wide updater requires elevation, it will be a separately signed, narrowly scoped executable with a versioned allowlisted protocol. It will not accept arbitrary command lines, filesystem paths outside approved install locations, proxy credentials, Android device data, or cloud refresh credentials.
