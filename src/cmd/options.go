package cmd

type globalOpts struct {
	downloadUrl     string
	followRedirects bool
}

var defaultGlobalOpts = globalOpts{"https://nodejs.org/dist", true}

type installOpts struct {
	skipChecksumValidation bool
	use                    bool
}

var defaultInstallOpts = installOpts{false, false}

type listOpts struct {
}

var defaultListOpts = listOpts{}

type uninstallOpts struct {
}

var defaultUninstallOpts = uninstallOpts{}

type useOpts struct {
}

var defaultUseOpts = useOpts{}
