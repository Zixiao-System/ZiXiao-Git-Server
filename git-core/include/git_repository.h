#ifndef GIT_REPOSITORY_H
#define GIT_REPOSITORY_H

#include <string>
#include <vector>

namespace GitCore {

class GitRepository {
public:
    GitRepository(const std::string& path);
    ~GitRepository();

    // Repository operations
    bool init(bool bare = true);
    bool exists() const;
    bool isValid() const;

    // Path operations
    std::string getPath() const;
    std::string getObjectsPath() const;
    std::string getRefsPath() const;
    std::string getHeadPath() const;

    // Reference operations
    bool createRef(const std::string& refName, const std::string& sha);
    std::string getRef(const std::string& refName) const;
    std::vector<std::string> listRefs() const;
    bool deleteRef(const std::string& refName);

    // Branch operations
    bool createBranch(const std::string& branchName, const std::string& sha);
    std::vector<std::string> listBranches() const;
    bool deleteBranch(const std::string& branchName);

    // Pack operations (for git protocol)
    bool receivePack(const std::string& packData);
    std::string uploadPack(const std::vector<std::string>& wants,
                          const std::vector<std::string>& haves);

private:
    std::string repoPath;
    bool initialized;

    bool createDirectory(const std::string& path);
    bool writeFile(const std::string& path, const std::string& content);
    std::string readFile(const std::string& path) const;
};

} // namespace GitCore

#endif // GIT_REPOSITORY_H
