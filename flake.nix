{
  description = " AWS Utilities for Cloud Engineers" ;

  inputs.nixpkgs.url = "nixpkgs/nixos-24.05";
  inputs.jsonify-aws-dotfiles.url = "github:wearetechnative/jsonify-aws-dotfiles";

  outputs = { self, nixpkgs, jsonify-aws-dotfiles }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      nixosModules.default = import ./module.nix self;

      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
          pkgJsonifyAwsDotfiles = jsonify-aws-dotfiles.packages.${system}.jsonify-aws-dotfiles;
        in
        {
          bmc = pkgs.callPackage ./package.nix { inherit pkgJsonifyAwsDotfiles;};
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.bmc);

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
          pkgJsonifyAwsDotfiles = jsonify-aws-dotfiles.packages.${system}.jsonify-aws-dotfiles;
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              awscli2
              aws-mfa
              granted
              jq
              dasel
              gum
              pkgJsonifyAwsDotfiles
            ];
          };
        });
    };
}
