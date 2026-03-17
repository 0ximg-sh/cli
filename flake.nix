{
  description = "0ximg CLI development environment and package";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, flake-utils, ... }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
    in
    flake-utils.lib.eachSystem systems (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };

        lib = pkgs.lib;
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            goreleaser
          ];

          shellHook = ''
            export CGO_ENABLED="0"
            echo "🚀 0ximg dev environment loaded!"
          '';
        };

        packages.default = pkgs.buildGoModule {
          pname = "0ximg";
          version = "0.1.0";
          src = ./.;

          vendorHash = "sha256-d0M5J72FKQLkZ97uftkGn53QCvT7FCyyNooyfbeapQk=";

          ldflags = [
            "-s"
            "-w"
          ];

          env.CGO_ENABLED = "0";

          postInstall = ''
            if [ -f "$out/bin/cli" ]; then
              mv "$out/bin/cli" "$out/bin/0ximg"
            fi
          '';

          meta = with lib; {
            description = "CLI for rendering code snippets with 0ximg";
            mainProgram = "0ximg";
            platforms = systems;
          };
        };
      });
}
