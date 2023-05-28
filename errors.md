# Common Errors

## Rollup Errors
`R0001: no minimal stake set, to fix set minimal stake for builder on payment contract`

This is an error you should only get on the startup of boost. To resolve this error, you need to set the minimal stake for the builder on the payment contract. This can be done by calling the `setMinimalStake` function on the payment contract. This value represents the amount of tokens that the searcher will need to stake in order to connect to your builder. When the payment is made, 80% of the funds staked, will be transfered to your builder address, and 20% will be transfered to the primev network.

**Error:** `failed to process next blocks`

This is that occurs while the rollup is processing data. If you encounter this issue, it may be a network issue with connecting to the RPC endpoint that hosts the rollup data.