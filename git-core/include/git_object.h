#ifndef GIT_OBJECT_H
#define GIT_OBJECT_H

#include <string>
#include <vector>
#include <cstdint>

namespace GitCore {

enum class GitObjectType {
    BLOB,
    TREE,
    COMMIT,
    TAG
};

class GitObject {
public:
    GitObject(GitObjectType type, const std::string& data);
    virtual ~GitObject();

    GitObjectType getType() const;
    std::string getSHA() const;
    std::string getData() const;
    size_t getSize() const;

    // Serialize object to git format
    std::string serialize() const;

    // Calculate SHA-1 hash
    static std::string calculateSHA(const std::string& content);

protected:
    GitObjectType type;
    std::string data;
    std::string sha;
};

class GitBlob : public GitObject {
public:
    explicit GitBlob(const std::string& content);
    std::string getContent() const;
};

class GitTreeEntry {
public:
    std::string mode;
    std::string name;
    std::string sha;

    GitTreeEntry(const std::string& mode, const std::string& name,
                 const std::string& sha);
};

class GitTree : public GitObject {
public:
    explicit GitTree(const std::vector<GitTreeEntry>& entries);

    void addEntry(const GitTreeEntry& entry);
    std::vector<GitTreeEntry> getEntries() const;

private:
    std::vector<GitTreeEntry> entries;
    std::string buildTreeData() const;
};

class GitCommit : public GitObject {
public:
    GitCommit(const std::string& treeSHA,
              const std::vector<std::string>& parentSHAs,
              const std::string& author,
              const std::string& committer,
              const std::string& message);

    std::string getTreeSHA() const;
    std::vector<std::string> getParentSHAs() const;
    std::string getMessage() const;

private:
    std::string treeSHA;
    std::vector<std::string> parentSHAs;
    std::string author;
    std::string committer;
    std::string message;

    std::string buildCommitData() const;
};

} // namespace GitCore

#endif // GIT_OBJECT_H
