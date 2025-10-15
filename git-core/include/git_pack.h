#ifndef GIT_PACK_H
#define GIT_PACK_H

#include <string>
#include <vector>
#include <cstdint>

namespace GitCore {

class GitPack {
public:
    GitPack();
    ~GitPack();

    // Pack file operations
    bool createPack(const std::vector<std::string>& objectSHAs,
                   std::string& packData);
    bool extractPack(const std::string& packData,
                    const std::string& objectsPath);

    // Pack index operations
    bool createIndex(const std::string& packPath,
                    const std::string& idxPath);

    struct PackObject {
        uint8_t type;
        uint64_t size;
        std::string data;
        std::string sha;
    };

    std::vector<PackObject> parsePackFile(const std::string& packData);

private:
    // Pack format constants
    static const uint32_t PACK_SIGNATURE = 0x5041434b; // 'PACK'
    static const uint32_t PACK_VERSION = 2;

    // Object types in pack
    enum PackObjectType {
        OBJ_COMMIT = 1,
        OBJ_TREE = 2,
        OBJ_BLOB = 3,
        OBJ_TAG = 4,
        OBJ_OFS_DELTA = 6,
        OBJ_REF_DELTA = 7
    };

    std::string compressData(const std::string& data);
    std::string decompressData(const std::string& compressed);

    uint64_t readVarint(const uint8_t* data, size_t& offset);
    void writeVarint(std::vector<uint8_t>& output, uint64_t value);
};

} // namespace GitCore

#endif // GIT_PACK_H
