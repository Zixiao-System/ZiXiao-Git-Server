#ifndef GIT_PROTOCOL_H
#define GIT_PROTOCOL_H

#include <string>
#include <vector>
#include <map>

namespace GitCore {

class GitProtocol {
public:
    GitProtocol();
    ~GitProtocol();

    // Pack line protocol
    static std::string pktLine(const std::string& data);
    static std::string flushPkt();
    static std::vector<std::string> parsePktLines(const std::string& input);

    // Git protocol commands
    struct RefAdvertisement {
        std::string sha;
        std::string refName;
        std::map<std::string, std::string> capabilities;
    };

    static std::string createRefAdvertisement(
        const std::vector<RefAdvertisement>& refs,
        const std::string& service);

    // git-receive-pack (push)
    struct PushRequest {
        std::vector<std::string> commands; // old-sha new-sha refname
        std::string packData;
    };

    static PushRequest parseReceivePack(const std::string& input);
    static std::string createReceivePackResponse(bool success,
                                                 const std::string& message);

    // git-upload-pack (pull/fetch)
    struct FetchRequest {
        std::vector<std::string> wants;
        std::vector<std::string> haves;
        int depth;
    };

    static FetchRequest parseUploadPack(const std::string& input);
    static std::string createUploadPackResponse(const std::string& packData);

private:
    static std::string encodeHex(const std::string& data);
    static std::string decodeHex(const std::string& hex);
};

} // namespace GitCore

#endif // GIT_PROTOCOL_H
