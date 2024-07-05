{
  description = "Flake for proprietary binary version of atlas cli";

  inputs = {
    nixpkgs.url = "flake:nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          atlas =
          let
            os = "linux";
            arch = "amd64";
            version = "v0.24.1-996d187-canary";
            ext = if os == "windows" then ".exe" else "";
          in
            pkgs.fetchurl {
              name = "atlas";
              url = "https://release.ariga.io/atlas/atlas-${os}-${arch}-${version}${ext}";
              sha256 = "sha256-wf0hTra1QZMdqGbd5qvujA+0fJ3hhPghSFLKCa1txyI=";

              recursiveHash = true;
              downloadToTemp = true;

              postFetch = ''
                set -ex
                pwd
                echo "$downloadedFile"
                mkdir -p "$out/bin"
                chmod +x "$downloadedFile"
                mv "$downloadedFile" "$out/bin/atlas"
              '';
            };
          default = self.packages.${system}.atlas;
        };
      }
    );
}
