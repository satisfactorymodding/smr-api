package tests

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/MarvinJWendt/testza"

	"github.com/satisfactorymodding/smr-api/validation"
)

func TestFilesToScan(t *testing.T) {
	ctx, _, stop := setup()
	defer stop()

	tests := []struct {
		name         string
		file         string
		modReference string
		want         []string
	}{
		{
			name:         "regular mod",
			file:         "SML-3.7.0.smod",
			modReference: "SML",
			want: []string{
				"Windows/Binaries/Win64/FactoryGame-SML-Win64-Shipping.dll",
				"WindowsServer/Binaries/Win64/UnrealServer-SML-Win64-Shipping.dll",
				"LinuxServer/Binaries/Linux/libUnrealServer-SML-Linux-Shipping.so",
			},
		},
		{
			name:         "non-UE files",
			file:         "ModWithFilesToScan.smod",
			modReference: "TestMod",
			want: []string{
				"extrafiles/helloworld",       // ELF executable
				"extrafiles/helloworld.exe",   // PE executable
				"extrafiles/libhelloworld.so", // ELF shared object
				"extrafiles/helloworld.dll",   // PE shared object

				// Not directly executable
				"extrafiles/helloworld.o", // ELF object file
				// "extrafiles/helloworld.obj", // COFF object file (gabriel-vasile/mimetype doesn't support this out of the box, but we could add it)
				// "extrafiles/libhelloworld.a", // ar archive containing ELF object files (is application/x-archive, we might want to check nested archives in general in the future)
				// "extrafiles/helloworld.lib", // ar archive containing ELF object files (is application/x-archive, we might want to check nested archives in general in the future)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := os.ReadFile(filepath.Join("testdata", tt.file))
			if err != nil {
				t.Fatalf("failed to read test mod file: %v", err)
			}

			got, err := validation.ExtractModInfo(ctx, bytes, false, tt.modReference)
			testza.AssertNoError(t, err)

			// Sort the slices before comparing, as the order may not be guaranteed
			slices.Sort(got.FilesToScan)
			slices.Sort(tt.want)

			testza.AssertEqual(t, tt.want, got.FilesToScan)
		})
	}
}
