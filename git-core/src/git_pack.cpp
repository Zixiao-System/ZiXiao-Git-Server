#include "git_pack.h"
#include <zlib.h>
#include <cstring>
#include <stdexcept>

namespace GitCore {

GitPack::GitPack() {
}

GitPack::~GitPack() {
}

bool GitPack::createPack(const std::vector<std::string>& objectSHAs,
                        std::string& packData) {
    // Simplified pack creation
    // In a real implementation, this would:
    // 1. Read all objects
    // 2. Compress them
    // 3. Write pack header, objects, and checksum

    std::vector<uint8_t> pack;

    // Write pack signature
    pack.push_back((PACK_SIGNATURE >> 24) & 0xFF);
    pack.push_back((PACK_SIGNATURE >> 16) & 0xFF);
    pack.push_back((PACK_SIGNATURE >> 8) & 0xFF);
    pack.push_back(PACK_SIGNATURE & 0xFF);

    // Write version
    pack.push_back((PACK_VERSION >> 24) & 0xFF);
    pack.push_back((PACK_VERSION >> 16) & 0xFF);
    pack.push_back((PACK_VERSION >> 8) & 0xFF);
    pack.push_back(PACK_VERSION & 0xFF);

    // Write object count
    uint32_t objCount = objectSHAs.size();
    pack.push_back((objCount >> 24) & 0xFF);
    pack.push_back((objCount >> 16) & 0xFF);
    pack.push_back((objCount >> 8) & 0xFF);
    pack.push_back(objCount & 0xFF);

    // Objects would be written here...

    packData = std::string(pack.begin(), pack.end());
    return true;
}

bool GitPack::extractPack(const std::string& packData,
                         const std::string& objectsPath) {
    // Simplified pack extraction
    if (packData.length() < 12) {
        return false;
    }

    // Verify pack signature
    uint32_t sig = (static_cast<uint8_t>(packData[0]) << 24) |
                   (static_cast<uint8_t>(packData[1]) << 16) |
                   (static_cast<uint8_t>(packData[2]) << 8) |
                   static_cast<uint8_t>(packData[3]);

    if (sig != PACK_SIGNATURE) {
        return false;
    }

    // Parse pack version
    uint32_t version = (static_cast<uint8_t>(packData[4]) << 24) |
                       (static_cast<uint8_t>(packData[5]) << 16) |
                       (static_cast<uint8_t>(packData[6]) << 8) |
                       static_cast<uint8_t>(packData[7]);

    if (version != PACK_VERSION) {
        return false;
    }

    // Get object count
    uint32_t objCount = (static_cast<uint8_t>(packData[8]) << 24) |
                        (static_cast<uint8_t>(packData[9]) << 16) |
                        (static_cast<uint8_t>(packData[10]) << 8) |
                        static_cast<uint8_t>(packData[11]);

    // Parse and extract objects...
    // In a real implementation, this would decompress and write objects

    return true;
}

bool GitPack::createIndex(const std::string& packPath,
                         const std::string& idxPath) {
    // Simplified index creation
    // In a real implementation, this would parse the pack file
    // and create an index for fast object lookup
    return true;
}

std::vector<GitPack::PackObject> GitPack::parsePackFile(const std::string& packData) {
    std::vector<PackObject> objects;

    if (packData.length() < 12) {
        return objects;
    }

    size_t offset = 12; // Skip header

    // Parse objects
    // This is a simplified version
    // Real implementation would parse each object header and data

    return objects;
}

std::string GitPack::compressData(const std::string& data) {
    z_stream zs;
    memset(&zs, 0, sizeof(zs));

    if (deflateInit(&zs, Z_DEFAULT_COMPRESSION) != Z_OK) {
        throw std::runtime_error("deflateInit failed");
    }

    zs.next_in = (Bytef*)data.data();
    zs.avail_in = data.size();

    int ret;
    char outbuffer[32768];
    std::string compressed;

    do {
        zs.next_out = reinterpret_cast<Bytef*>(outbuffer);
        zs.avail_out = sizeof(outbuffer);

        ret = deflate(&zs, Z_FINISH);

        if (compressed.size() < zs.total_out) {
            compressed.append(outbuffer, zs.total_out - compressed.size());
        }
    } while (ret == Z_OK);

    deflateEnd(&zs);

    if (ret != Z_STREAM_END) {
        throw std::runtime_error("deflate failed");
    }

    return compressed;
}

std::string GitPack::decompressData(const std::string& compressed) {
    z_stream zs;
    memset(&zs, 0, sizeof(zs));

    if (inflateInit(&zs) != Z_OK) {
        throw std::runtime_error("inflateInit failed");
    }

    zs.next_in = (Bytef*)compressed.data();
    zs.avail_in = compressed.size();

    int ret;
    char outbuffer[32768];
    std::string decompressed;

    do {
        zs.next_out = reinterpret_cast<Bytef*>(outbuffer);
        zs.avail_out = sizeof(outbuffer);

        ret = inflate(&zs, 0);

        if (decompressed.size() < zs.total_out) {
            decompressed.append(outbuffer, zs.total_out - decompressed.size());
        }
    } while (ret == Z_OK);

    inflateEnd(&zs);

    if (ret != Z_STREAM_END) {
        throw std::runtime_error("inflate failed");
    }

    return decompressed;
}

uint64_t GitPack::readVarint(const uint8_t* data, size_t& offset) {
    uint64_t value = 0;
    uint8_t byte;
    int shift = 0;

    do {
        byte = data[offset++];
        value |= ((uint64_t)(byte & 0x7F)) << shift;
        shift += 7;
    } while (byte & 0x80);

    return value;
}

void GitPack::writeVarint(std::vector<uint8_t>& output, uint64_t value) {
    while (value > 0x7F) {
        output.push_back((value & 0x7F) | 0x80);
        value >>= 7;
    }
    output.push_back(value & 0x7F);
}

} // namespace GitCore
