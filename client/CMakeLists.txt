cmake_minimum_required(VERSION 3.0)
project(CloudBoard)

find_package(PkgConfig REQUIRED)
pkg_check_modules(GTK3 REQUIRED gtk+-3.0)

set(SOURCES clipboard-daemon.cpp)

include_directories(${GTK3_INCLUDE_DIRS})
link_directories(${GTK3_INCLUDE_DIRS})
add_definitions(${GTK3_INCLUDE_DIRS})

add_executable(cboardd ${SOURCES})
target_link_libraries(cboardd ${GTK3_LIBRARIES})
