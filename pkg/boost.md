
# Boost Service
The Boost service is responsible for processing block data into metadata and making it available for searcher consumption.

## Type Definitions
**Boost**: This is an interface that defines two functions - SubmitBlock and GetWorkChannel.

**DefaultBoost**: This is the default implementation of the Boost interface, containing a configuration object and a channel to push metadata.

**Transaction**: This struct contains information related to a transaction, including the count and priority fee details.

**Metadata**: This struct contains details about the block, including the builder, number, block hash, timestamp, base fee, and transactions.

## Main Functions
**NewBoost**: This function is the constructor for the DefaultBoost struct. It takes a configuration object, validates it, and creates a new DefaultBoost object.

**SubmitBlock**: This function processes a SubmitBlockRequest from your Execution client. It extracts the block and transaction metadata and pushes the metadata into a channel for consumption by searchers. If any error occurs during the processing, it recovers from the panic and logs an error message.

## Error Definitions
**ErrBlockUnprocessable**: This error is thrown when the SubmitBlock function encounters an unprocessable block.

# Processing Flow
1. A block submission request (SubmitBlockRequest) is received.

2. The block details such as block hash, block number, builder pubkey, transaction count, timestamp, and base fee are extracted and saved in a Metadata object.

3. If there are transactions in the block, the function calculates the minimum and maximum priority fees among all the transactions.

4. The Metadata object is then pushed into a channel for publishing to connected searchers via a websocket.

5. Finally, the function logs the details of the processed block.
