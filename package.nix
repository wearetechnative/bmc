{ lib
, buildGoModule
}:

buildGoModule rec {
  pname = "bmc";
  version = lib.fileContents ./VERSION-bmc;

  src = ./.;

  vendorHash = "sha256-1hHNNsJuc93f9WL/aDpBBnuz6+9b7K9fJiXwuuQIepg=";

  ldflags = [ "-s" "-w" "-X github.com/wearetechnative/bmc/cmd.Version=${version}" ];

  meta = with lib; {
    description = "Bill McCloud's AWS toolbox — profile selection, EC2/ECS operations, console access";
    homepage = "https://github.com/wearetechnative/bmc";
    license = licenses.asl20;
    mainProgram = "bmc";
  };
}
