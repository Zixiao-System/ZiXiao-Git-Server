#include "git_protocol.h"
#include <sstream>
#include <iomanip>

namespace GitCore {

GitProtocol::GitProtocol() {
}

GitProtocol::~GitProtocol() {
}

std::string GitProtocol::pktLine(const std::string& data) {
    if (data.empty()) {
        return "0000";
    }

    // Length includes 4 bytes for length itself
    int len = data.length() + 4;
    if (len > 65524) { // Max pkt-line length
        return "0000";
    }

    std::ostringstream oss;
    oss << std::hex << std::setw(4) << std::setfill('0') << len;
    oss << data;
    return oss.str();
}

std::string GitProtocol::flushPkt() {
    return "0000";
}

std::vector<std::string> GitProtocol::parsePktLines(const std::string& input) {
    std::vector<std::string> lines;
    size_t pos = 0;

    while (pos < input.length()) {
        if (pos + 4 > input.length()) {
            break;
        }

        std::string lenStr = input.substr(pos, 4);
        int len = std::stoi(lenStr, nullptr, 16);

        if (len == 0) {
            // Flush packet
            pos += 4;
            lines.push_back("");
            continue;
        }

        if (len < 4 || pos + len > input.length()) {
            break;
        }

        std::string line = input.substr(pos + 4, len - 4);
        lines.push_back(line);
        pos += len;
    }

    return lines;
}

std::string GitProtocol::createRefAdvertisement(
    const std::vector<RefAdvertisement>& refs,
    const std::string& service) {

    std::ostringstream oss;

    // Service announcement
    oss << pktLine("# service=" + service + "\n");
    oss << flushPkt();

    if (refs.empty()) {
        // No refs, advertise capabilities only
        oss << pktLine("0000000000000000000000000000000000000000 "
                      "capabilities^{}\0report-status delete-refs "
                      "side-band-64k quiet atomic\n");
    } else {
        // First ref includes capabilities
        bool first = true;
        for (const auto& ref : refs) {
            std::string line = ref.sha + " " + ref.refName;

            if (first) {
                line += "\0report-status delete-refs side-band-64k quiet atomic";
                first = false;
            }

            line += "\n";
            oss << pktLine(line);
        }
    }

    oss << flushPkt();
    return oss.str();
}

GitProtocol::PushRequest GitProtocol::parseReceivePack(const std::string& input) {
    PushRequest request;
    std::vector<std::string> lines = parsePktLines(input);

    bool commandsDone = false;
    std::string packData;

    for (const auto& line : lines) {
        if (line.empty()) {
            commandsDone = true;
            continue;
        }

        if (!commandsDone) {
            // Parse command: old-sha new-sha refname
            request.commands.push_back(line);
        } else {
            // Pack data
            packData += line;
        }
    }

    request.packData = packData;
    return request;
}

std::string GitProtocol::createReceivePackResponse(bool success,
                                                   const std::string& message) {
    std::ostringstream oss;

    if (success) {
        oss << pktLine("unpack ok\n");
    } else {
        oss << pktLine("unpack " + message + "\n");
    }

    oss << flushPkt();
    return oss.str();
}

GitProtocol::FetchRequest GitProtocol::parseUploadPack(const std::string& input) {
    FetchRequest request;
    request.depth = 0;

    std::vector<std::string> lines = parsePktLines(input);

    for (const auto& line : lines) {
        if (line.empty()) {
            continue;
        }

        if (line.substr(0, 5) == "want ") {
            std::string sha = line.substr(5, 40);
            request.wants.push_back(sha);
        } else if (line.substr(0, 5) == "have ") {
            std::string sha = line.substr(5, 40);
            request.haves.push_back(sha);
        } else if (line.substr(0, 6) == "depth ") {
            request.depth = std::stoi(line.substr(6));
        }
    }

    return request;
}

std::string GitProtocol::createUploadPackResponse(const std::string& packData) {
    std::ostringstream oss;

    // NAK response
    oss << pktLine("NAK\n");

    // Pack data
    oss << pktLine(packData);
    oss << flushPkt();

    return oss.str();
}

std::string GitProtocol::encodeHex(const std::string& data) {
    std::ostringstream oss;
    for (unsigned char c : data) {
        oss << std::hex << std::setw(2) << std::setfill('0')
            << static_cast<int>(c);
    }
    return oss.str();
}

std::string GitProtocol::decodeHex(const std::string& hex) {
    std::string result;
    for (size_t i = 0; i < hex.length(); i += 2) {
        std::string byteStr = hex.substr(i, 2);
        unsigned char byte = static_cast<unsigned char>(
            std::stoi(byteStr, nullptr, 16));
        result += byte;
    }
    return result;
}

} // namespace GitCore
