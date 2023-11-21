O formato de arquivo compactado é o seguinte: (big-endian)

HEADER
	- START_FLAG				        2 bytes (uint16)
	- SRC_FILENAME_LEN			        2 bytes (uint16)
	- BYTE SIZE BEFORE COMPRESSION		4 bytes (uint32)
	- BYTE SIZE AFTER COMPRESSION		4 bytes (uint32)
	- SRC_FILENAME				        n bytes

DATA
	- HUFFMAN TABLE
		-- HUFFMAN TABLE SIZE 		    4 bytes (uint32)
		-- HUFFMAN TABLE DATA
	- COMPRESSED DATA
		-- VALID BIT LEN		        4 bytes (uint32) + 1 bytes = 5 bytes
		-- COMPRESSED BIT

TAIL
	- CRC32 CHECKSUM	  	            4 bytes (uint32)
	- END_FLAG			                2 bytes (uint16)



O formato de serialização é o seguinte (big-endian: bits altos são colocados no endereço baixo e os bits baixos são colocados no endereço alto)

START_FLAG			        4 bytes
NUMBER OF TABLE ITEMS		4 bytes (uint32)
TABLE_ITEM_1(BYTE+CODE)		1+4=5 bytes
TABLE_ITEM_2(BYTE+CODE)		1+4=5 bytes
...
TABLE_ITEM_N(BYTE+CODE)		1+4=5 bytes
CRC32				        4 bytes
END_FLAG			        4 bytes