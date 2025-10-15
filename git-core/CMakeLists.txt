cmake_minimum_required(VERSION 3.15)
project(ZiXiaoGitCore VERSION 1.0.0 LANGUAGES CXX)

# C++ Standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# Build type
if(NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE Release)
endif()

# Compiler flags
if(MSVC)
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} /W4 /EHsc")
    set(CMAKE_CXX_FLAGS_RELEASE "/O2 /MD")
    add_definitions(-D_CRT_SECURE_NO_WARNINGS)
else()
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -Wextra")
    set(CMAKE_CXX_FLAGS_DEBUG "-g -O0")
    set(CMAKE_CXX_FLAGS_RELEASE "-O2")
endif()

# Find required packages
find_package(OpenSSL REQUIRED)
find_package(ZLIB REQUIRED)

# macOS specific settings
if(APPLE)
    # Homebrew OpenSSL paths
    if(EXISTS "/opt/homebrew/opt/openssl")
        set(OPENSSL_ROOT_DIR "/opt/homebrew/opt/openssl")
    elseif(EXISTS "/usr/local/opt/openssl")
        set(OPENSSL_ROOT_DIR "/usr/local/opt/openssl")
    endif()

    include_directories(${OPENSSL_ROOT_DIR}/include)
    link_directories(${OPENSSL_ROOT_DIR}/lib)
endif()

# Windows specific settings
if(WIN32)
    # Use vcpkg if available
    if(DEFINED ENV{VCPKG_ROOT})
        include_directories($ENV{VCPKG_ROOT}/installed/x64-windows/include)
        link_directories($ENV{VCPKG_ROOT}/installed/x64-windows/lib)
    endif()
endif()

# Include directories
include_directories(
    ${CMAKE_CURRENT_SOURCE_DIR}/include
    ${OPENSSL_INCLUDE_DIR}
)

# Source files
set(SOURCES
    src/git_repository.cpp
    src/git_object.cpp
    src/git_protocol.cpp
    src/git_pack.cpp
    src/git_c_api.cpp
)

# Header files
set(HEADERS
    include/git_repository.h
    include/git_object.h
    include/git_protocol.h
    include/git_pack.h
    include/git_c_api.h
)

# Create shared library
add_library(gitcore SHARED ${SOURCES} ${HEADERS})

# Link libraries
target_link_libraries(gitcore
    OpenSSL::SSL
    OpenSSL::Crypto
    ZLIB::ZLIB
)

# Windows specific libraries
if(WIN32)
    target_link_libraries(gitcore ws2_32 crypt32)
endif()

# Set library properties
set_target_properties(gitcore PROPERTIES
    VERSION ${PROJECT_VERSION}
    SOVERSION 1
    PUBLIC_HEADER "${HEADERS}"
    WINDOWS_EXPORT_ALL_SYMBOLS ON
)

# Installation
install(TARGETS gitcore
    LIBRARY DESTINATION lib
    RUNTIME DESTINATION bin
    ARCHIVE DESTINATION lib
    PUBLIC_HEADER DESTINATION include/gitcore
)

# Output directory
set_target_properties(gitcore PROPERTIES
    LIBRARY_OUTPUT_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/lib
    RUNTIME_OUTPUT_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/lib
    ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/lib
)

# Print configuration
message(STATUS "Build type: ${CMAKE_BUILD_TYPE}")
message(STATUS "C++ compiler: ${CMAKE_CXX_COMPILER}")
message(STATUS "OpenSSL version: ${OPENSSL_VERSION}")
message(STATUS "OpenSSL include: ${OPENSSL_INCLUDE_DIR}")
message(STATUS "ZLIB version: ${ZLIB_VERSION_STRING}")
message(STATUS "Platform: ${CMAKE_SYSTEM_NAME}")

