## Solution

### Representation

Any drum binary would contain following information in the same order

1. 6 byte - SPLICE header (Bigendian)
2. 7 dummy bytes
3. 1 byte - size
4. 11 byte - version
5. 21 dummy bytes
6. float32 - Tempo (Little Endian)
7. 1 dummy byte
7. Tracks
   - 1 byte - track id
   - 2 dummy bytes
   - 1 byte length of the track title
   - length's name
   - 16 byte for the 16 states
