cmake_minimum_required(VERSION 3.12)
project(MapleSAT)

# Set the executable name
set(EXEC maplesat)

# Add source files
file(GLOB CSRCS "${CMAKE_CURRENT_SOURCE_DIR}/**/*.cc")
list(REMOVE_ITEM CSRCS ${CMAKE_CURRENT_SOURCE_DIR}/core/Main.cc ${CMAKE_CURRENT_SOURCE_DIR}/test/test.cc)
# include(CMakePrintHelpers)
# cmake_print_variables(CSRCS)

# Compiler options
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -Wno-parentheses -O3 -g -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS")

# Include directories
include_directories(${CMAKE_CURRENT_SOURCE_DIR})

# Build executable
add_executable(${EXEC} ${CSRCS})
target_link_libraries(${EXEC} ntl z)
