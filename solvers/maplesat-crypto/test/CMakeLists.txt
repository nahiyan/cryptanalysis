cmake_minimum_required (VERSION 2.8.11)

include_directories("../")

add_executable(test test.cc ../core/Crypto.cc ../core/Solver.cc)
target_link_libraries(test ntl)
