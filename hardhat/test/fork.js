const { hre } = require("hardhat");

const forkFrom = async (blockNumber) => {
  if (!hre.config.networks.hardhat.forking) {
    throw new Error(
      `Forking misconfigured for "hardhat" network in hardhat.config.ts`
    );
  }

  await hre.network.provider.request({
    method: "hardhat_reset",
    params: [
      {
        forking: {
          jsonRpcUrl: hre.config.networks.hardhat.forking.url,
          blockNumber: blockNumber,
        },
      },
    ],
  });
};