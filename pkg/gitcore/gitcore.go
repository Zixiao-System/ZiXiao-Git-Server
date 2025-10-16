package gitcore

/*
#cgo CXXFLAGS: -std=c++17 -I${SRCDIR}/../git-core/include
#cgo darwin LDFLAGS: -L${SRCDIR}/../git-core/lib -lgitcore -lstdc++ -lz -L/opt/homebrew/opt/openssl/lib -L/usr/local/opt/openssl/lib -lcrypto -lssl
#cgo linux LDFLAGS: -L${SRCDIR}/../git-core/lib -lgitcore -lstdc++ -lz -lcrypto -lssl
#include "git_c_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

// Repository wraps the C++ GitRepository
type Repository struct {
	ptr unsafe.Pointer
}

// NewRepository creates a new repository instance
func NewRepository(path string) *Repository {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ptr := C.git_repository_new(cPath)
	return &Repository{ptr: ptr}
}

// Free releases the repository resources
func (r *Repository) Free() {
	if r.ptr != nil {
		C.git_repository_free(r.ptr)
		r.ptr = nil
	}
}

// Init initializes a new git repository
func (r *Repository) Init(bare bool) error {
	bareInt := 0
	if bare {
		bareInt = 1
	}

	result := C.git_repository_init(r.ptr, C.int(bareInt))
	if result == 0 {
		return errors.New("failed to initialize repository")
	}
	return nil
}

// Exists checks if the repository exists
func (r *Repository) Exists() bool {
	result := C.git_repository_exists(r.ptr)
	return result != 0
}

// IsValid checks if the repository is valid
func (r *Repository) IsValid() bool {
	result := C.git_repository_is_valid(r.ptr)
	return result != 0
}

// CreateRef creates a new reference
func (r *Repository) CreateRef(refName, sha string) error {
	cRefName := C.CString(refName)
	cSha := C.CString(sha)
	defer C.free(unsafe.Pointer(cRefName))
	defer C.free(unsafe.Pointer(cSha))

	result := C.git_repository_create_ref(r.ptr, cRefName, cSha)
	if result == 0 {
		return errors.New("failed to create reference")
	}
	return nil
}

// GetRef retrieves a reference value
func (r *Repository) GetRef(refName string) (string, error) {
	cRefName := C.CString(refName)
	defer C.free(unsafe.Pointer(cRefName))

	cResult := C.git_repository_get_ref(r.ptr, cRefName)
	if cResult == nil {
		return "", errors.New("reference not found")
	}
	defer C.git_free_string(cResult)

	return C.GoString(cResult), nil
}

// ListRefs lists all references
func (r *Repository) ListRefs() ([]string, error) {
	var count C.int
	cRefs := C.git_repository_list_refs(r.ptr, &count)
	if cRefs == nil {
		return []string{}, nil
	}
	defer C.git_free_string_array(cRefs, count)

	refs := make([]string, int(count))
	refSlice := (*[1 << 28]*C.char)(unsafe.Pointer(cRefs))[:count:count]
	for i, cRef := range refSlice {
		refs[i] = C.GoString(cRef)
	}

	return refs, nil
}

// DeleteRef deletes a reference
func (r *Repository) DeleteRef(refName string) error {
	cRefName := C.CString(refName)
	defer C.free(unsafe.Pointer(cRefName))

	result := C.git_repository_delete_ref(r.ptr, cRefName)
	if result == 0 {
		return errors.New("failed to delete reference")
	}
	return nil
}

// CreateBranch creates a new branch
func (r *Repository) CreateBranch(branchName, sha string) error {
	cBranchName := C.CString(branchName)
	cSha := C.CString(sha)
	defer C.free(unsafe.Pointer(cBranchName))
	defer C.free(unsafe.Pointer(cSha))

	result := C.git_repository_create_branch(r.ptr, cBranchName, cSha)
	if result == 0 {
		return errors.New("failed to create branch")
	}
	return nil
}

// ListBranches lists all branches
func (r *Repository) ListBranches() ([]string, error) {
	var count C.int
	cBranches := C.git_repository_list_branches(r.ptr, &count)
	if cBranches == nil {
		return []string{}, nil
	}
	defer C.git_free_string_array(cBranches, count)

	branches := make([]string, int(count))
	branchSlice := (*[1 << 28]*C.char)(unsafe.Pointer(cBranches))[:count:count]
	for i, cBranch := range branchSlice {
		branches[i] = C.GoString(cBranch)
	}

	return branches, nil
}

// DeleteBranch deletes a branch
func (r *Repository) DeleteBranch(branchName string) error {
	cBranchName := C.CString(branchName)
	defer C.free(unsafe.Pointer(cBranchName))

	result := C.git_repository_delete_branch(r.ptr, cBranchName)
	if result == 0 {
		return errors.New("failed to delete branch")
	}
	return nil
}

// ReceivePack processes a git push pack
func (r *Repository) ReceivePack(packData []byte) error {
	cPackData := C.CBytes(packData)
	defer C.free(cPackData)

	result := C.git_repository_receive_pack(r.ptr, (*C.char)(cPackData), C.int(len(packData)))
	if result == 0 {
		return errors.New("failed to receive pack")
	}
	return nil
}

// UploadPack generates a pack for git pull/fetch
func (r *Repository) UploadPack(wants, haves []string) ([]byte, error) {
	// Convert wants to C array
	cWants := make([]*C.char, len(wants))
	for i, want := range wants {
		cWants[i] = C.CString(want)
		defer C.free(unsafe.Pointer(cWants[i]))
	}

	// Convert haves to C array
	cHaves := make([]*C.char, len(haves))
	for i, have := range haves {
		cHaves[i] = C.CString(have)
		defer C.free(unsafe.Pointer(cHaves[i]))
	}

	var outLen C.int
	var cWantsPtr **C.char
	var cHavesPtr **C.char

	if len(wants) > 0 {
		cWantsPtr = &cWants[0]
	}
	if len(haves) > 0 {
		cHavesPtr = &cHaves[0]
	}

	cResult := C.git_repository_upload_pack(r.ptr, cWantsPtr, C.int(len(wants)),
		cHavesPtr, C.int(len(haves)), &outLen)
	if cResult == nil {
		return nil, errors.New("failed to upload pack")
	}
	defer C.git_free_string(cResult)

	return C.GoBytes(unsafe.Pointer(cResult), outLen), nil
}

// Protocol functions

// PktLine creates a git protocol packet line
func PktLine(data string) string {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	cResult := C.git_protocol_pkt_line(cData)
	defer C.git_free_string(cResult)

	return C.GoString(cResult)
}

// FlushPkt creates a flush packet
func FlushPkt() string {
	cResult := C.git_protocol_flush_pkt()
	defer C.git_free_string(cResult)

	return C.GoString(cResult)
}

// CreateRefAdvertisement creates a reference advertisement for git protocol
func CreateRefAdvertisement(refs map[string]string, service string) ([]byte, error) {
	if len(refs) == 0 {
		return nil, errors.New("no refs to advertise")
	}

	// Convert map to C arrays
	cRefs := make([]*C.char, 0, len(refs))
	cShas := make([]*C.char, 0, len(refs))

	for ref, sha := range refs {
		cRefs = append(cRefs, C.CString(ref))
		cShas = append(cShas, C.CString(sha))
	}

	defer func() {
		for _, cRef := range cRefs {
			C.free(unsafe.Pointer(cRef))
		}
		for _, cSha := range cShas {
			C.free(unsafe.Pointer(cSha))
		}
	}()

	cService := C.CString(service)
	defer C.free(unsafe.Pointer(cService))

	var outLen C.int
	cResult := C.git_protocol_create_ref_advertisement(&cRefs[0], &cShas[0],
		C.int(len(refs)), cService, &outLen)
	if cResult == nil {
		return nil, errors.New("failed to create ref advertisement")
	}
	defer C.git_free_string(cResult)

	return C.GoBytes(unsafe.Pointer(cResult), outLen), nil
}
