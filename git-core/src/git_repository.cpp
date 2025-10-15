#include "git_repository.h"
#include <sys/stat.h>
#include <fstream>
#include <sstream>
#include <algorithm>
#include <filesystem>

namespace fs = std::filesystem;

namespace GitCore {

GitRepository::GitRepository(const std::string& path)
    : repoPath(path), initialized(false) {
}

GitRepository::~GitRepository() {
}

bool GitRepository::init(bool bare) {
    if (exists()) {
        return false;
    }

    // Create repository directory structure
    if (!createDirectory(repoPath)) {
        return false;
    }

    // Create git directory structure
    std::string gitDir = bare ? repoPath : repoPath + "/.git";

    if (!createDirectory(gitDir + "/objects")) return false;
    if (!createDirectory(gitDir + "/objects/pack")) return false;
    if (!createDirectory(gitDir + "/objects/info")) return false;
    if (!createDirectory(gitDir + "/refs")) return false;
    if (!createDirectory(gitDir + "/refs/heads")) return false;
    if (!createDirectory(gitDir + "/refs/tags")) return false;

    // Create HEAD file
    std::string headContent = "ref: refs/heads/main\n";
    if (!writeFile(gitDir + "/HEAD", headContent)) {
        return false;
    }

    // Create config file
    std::string configContent = "[core]\n\trepositoryformatversion = 0\n";
    if (bare) {
        configContent += "\tbare = true\n";
    } else {
        configContent += "\tfilemode = true\n";
    }
    if (!writeFile(gitDir + "/config", configContent)) {
        return false;
    }

    // Create description file
    if (!writeFile(gitDir + "/description",
                  "Unnamed repository; edit this file to name it.\n")) {
        return false;
    }

    initialized = true;
    return true;
}

bool GitRepository::exists() const {
    return fs::exists(repoPath);
}

bool GitRepository::isValid() const {
    if (!exists()) return false;

    std::string gitDir = repoPath;
    if (fs::exists(repoPath + "/.git")) {
        gitDir = repoPath + "/.git";
    }

    return fs::exists(gitDir + "/objects") &&
           fs::exists(gitDir + "/refs") &&
           fs::exists(gitDir + "/HEAD");
}

std::string GitRepository::getPath() const {
    return repoPath;
}

std::string GitRepository::getObjectsPath() const {
    std::string gitDir = repoPath;
    if (fs::exists(repoPath + "/.git")) {
        gitDir = repoPath + "/.git";
    }
    return gitDir + "/objects";
}

std::string GitRepository::getRefsPath() const {
    std::string gitDir = repoPath;
    if (fs::exists(repoPath + "/.git")) {
        gitDir = repoPath + "/.git";
    }
    return gitDir + "/refs";
}

std::string GitRepository::getHeadPath() const {
    std::string gitDir = repoPath;
    if (fs::exists(repoPath + "/.git")) {
        gitDir = repoPath + "/.git";
    }
    return gitDir + "/HEAD";
}

bool GitRepository::createRef(const std::string& refName, const std::string& sha) {
    std::string refPath = getRefsPath() + "/" + refName;

    // Create parent directories if needed
    size_t pos = refPath.find_last_of('/');
    if (pos != std::string::npos) {
        std::string parentDir = refPath.substr(0, pos);
        if (!fs::exists(parentDir)) {
            fs::create_directories(parentDir);
        }
    }

    return writeFile(refPath, sha + "\n");
}

std::string GitRepository::getRef(const std::string& refName) const {
    std::string refPath = getRefsPath() + "/" + refName;
    std::string content = readFile(refPath);

    if (content.empty()) {
        return "";
    }

    // Remove trailing newline
    if (!content.empty() && content.back() == '\n') {
        content.pop_back();
    }

    return content;
}

std::vector<std::string> GitRepository::listRefs() const {
    std::vector<std::string> refs;
    std::string refsPath = getRefsPath();

    for (const auto& entry : fs::recursive_directory_iterator(refsPath)) {
        if (entry.is_regular_file()) {
            std::string refPath = entry.path().string();
            // Get relative path from refs directory
            std::string refName = refPath.substr(refsPath.length() + 1);
            refs.push_back(refName);
        }
    }

    return refs;
}

bool GitRepository::deleteRef(const std::string& refName) {
    std::string refPath = getRefsPath() + "/" + refName;
    if (!fs::exists(refPath)) {
        return false;
    }
    return fs::remove(refPath);
}

bool GitRepository::createBranch(const std::string& branchName, const std::string& sha) {
    return createRef("heads/" + branchName, sha);
}

std::vector<std::string> GitRepository::listBranches() const {
    std::vector<std::string> branches;
    std::vector<std::string> allRefs = listRefs();

    for (const auto& ref : allRefs) {
        if (ref.find("heads/") == 0) {
            branches.push_back(ref.substr(6)); // Remove "heads/" prefix
        }
    }

    return branches;
}

bool GitRepository::deleteBranch(const std::string& branchName) {
    return deleteRef("heads/" + branchName);
}

bool GitRepository::receivePack(const std::string& packData) {
    // Simplified pack receiving
    // In a real implementation, this would parse and extract the pack file
    std::string packPath = getObjectsPath() + "/pack/pack-" +
                          std::to_string(std::time(nullptr)) + ".pack";
    return writeFile(packPath, packData);
}

std::string GitRepository::uploadPack(const std::vector<std::string>& wants,
                                     const std::vector<std::string>& haves) {
    // Simplified pack generation
    // In a real implementation, this would generate a pack file with requested objects
    return "PACK data placeholder";
}

bool GitRepository::createDirectory(const std::string& path) {
    try {
        return fs::create_directories(path);
    } catch (const std::exception& e) {
        return false;
    }
}

bool GitRepository::writeFile(const std::string& path, const std::string& content) {
    std::ofstream file(path, std::ios::binary);
    if (!file) {
        return false;
    }
    file << content;
    return file.good();
}

std::string GitRepository::readFile(const std::string& path) const {
    std::ifstream file(path, std::ios::binary);
    if (!file) {
        return "";
    }

    std::stringstream buffer;
    buffer << file.rdbuf();
    return buffer.str();
}

} // namespace GitCore
