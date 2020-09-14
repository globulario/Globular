TEMPLATE = app
CONFIG += console c++
CONFIG -= app_bundle
CONFIG -= qt
CONFIG += c++17

SOURCES += \
        ../../api/cpp/GlobularClient/globularclient.cpp \
        ../../ressource/cpp/GlobularRessourceClient/globularressourceclient.cpp \
        ../../ressource/cpp/GlobularRessourceClient/ressource/ressource.grpc.pb.cc \
        ../../ressource/cpp/GlobularRessourceClient/ressource/ressource.pb.cc \
        ../GlobularServer/globularserver.cpp \
        echo/echopb/echo.grpc.pb.cc \
        echo/echopb/echo.pb.cc \
        echoserviceimpl.cpp \
        main.cpp

HEADERS += \
    ../../api/cpp/GlobularClient/globularclient.h \
    ../../ressource/cpp/GlobularRessourceClient/globularressourceclient.h \
    ../../ressource/cpp/GlobularRessourceClient/ressource/ressource.grpc.pb.h \
    ../../ressource/cpp/GlobularRessourceClient/ressource/ressource.pb.h \
    ../GlobularServer/globularserver.h \
    echo/echopb/echo.grpc.pb.h \
    echo/echopb/echo.pb.h \
    echoserviceimpl.h

INCLUDEPATH +=  $$PWD/../../ressource/cpp/GlobularRessourceClient $$PWD/../../api/cpp/GlobularClient ../../cpp
INCLUDEPATH +=  $$PWD/../GlobularServer $$PWD/../../api/cpp/GlobularClient  $$PWD/../../ressource/cpp/GlobularRessourceClient $$PWD/../../ressource/cpp/GlobularRessourceClient

#here I will make use of pkg-config to get the list of dependencie of each libraries.
unix: LIBS += `pkg-config --libs grpc++ protobuf`

# Set the pkconfig.
win32: LIBS += -LC:/msys64/mingw64/lib -lgrpc++ -labsl_raw_hash_set -labsl_hashtablez_sampler -labsl_exponential_biased -labsl_hash -labsl_bad_variant_access -labsl_city -labsl_status -labsl_cord -labsl_bad_optional_access -labsl_str_format_internal -labsl_synchronization -labsl_graphcycles_internal -labsl_symbolize -labsl_demangle_internal -labsl_stacktrace -labsl_debugging_internal -labsl_malloc_internal -labsl_time -labsl_time_zone -labsl_civil_time -labsl_strings -labsl_strings_internal -labsl_throw_delegate -labsl_int128 -labsl_base -labsl_spinlock_wait -labsl_raw_logging_internal -labsl_log_severity -labsl_dynamic_annotations -lgrpc -laddress_sorting -lre2 -lupb -lcares -lz -labsl_raw_hash_set -labsl_hashtablez_sampler -labsl_exponential_biased -labsl_hash -labsl_bad_variant_access -labsl_city -labsl_status -labsl_cord -labsl_bad_optional_access -labsl_str_format_internal -labsl_synchronization -labsl_graphcycles_internal -labsl_symbolize -labsl_demangle_internal -labsl_stacktrace -labsl_debugging_internal -labsl_malloc_internal -labsl_time -labsl_time_zone -labsl_civil_time -labsl_strings -labsl_strings_internal -labsl_throw_delegate -labsl_int128 -labsl_base -labsl_spinlock_wait -labsl_raw_logging_internal -labsl_log_severity -labsl_dynamic_annotations -lgpr -labsl_str_format_internal -labsl_synchronization -labsl_graphcycles_internal -labsl_symbolize -labsl_demangle_internal -labsl_stacktrace -labsl_debugging_internal -labsl_malloc_internal -labsl_time -labsl_time_zone -labsl_civil_time -labsl_strings -labsl_strings_internal -labsl_throw_delegate -labsl_int128 -labsl_base -labsl_spinlock_wait -labsl_raw_logging_internal -labsl_log_severity -labsl_dynamic_annotations -lssl -lcrypto -lws2_32 -lgdi32 -lcrypt32  -limagehlp -lprotobuf
