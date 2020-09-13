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

INCLUDEPATH += C:\Users\mm006819\grpc\third_party\protobuf\src C:\Users\mm006819\grpc\include $$PWD/../../ressource/cpp/GlobularRessourceClient $$PWD/../../api/cpp/GlobularClient ../../cpp
INCLUDEPATH +=  $$PWD/../GlobularServer $$PWD/../../api/cpp/GlobularClient  $$PWD/../../ressource/cpp/GlobularRessourceClient $$PWD/../../ressource/cpp/GlobularRessourceClient

#grpc stuff...
#unix: LIBS += -labsl_bad_optional_access -labsl_bad_variant_access -labsl_base -labsl_city -labsl_civil_time -labsl_cord -labsl_debugging_internal -labsl_demangle_internal
#unix: LIBS += -labsl_dynamic_annotations -labsl_exponential_biased -labsl_graphcycles_internal -labsl_hash -labsl_hashtablez_sampler -labsl_int128 -labsl_log_severity
#unix: LIBS += -labsl_malloc_internal -labsl_raw_hash_set -labsl_raw_logging_internal -labsl_spinlock_wait -labsl_stacktrace -labsl_status -labsl_str_format_internal
#unix: LIBS += -labsl_strings -labsl_strings_internal -labsl_symbolize -labsl_synchronization -labsl_throw_delegate -labsl_time -labsl_time_zone
#unix: LIBS += -L/usr/local/lib -lgpr -lgrpc++ -lgrpc -lgrpc++_alts -lgrpc++_error_details -lgrpc_plugin_support -lgrpcpp_channelz
#unix: LIBS +=  -lgrpc++_reflection -lprotobuf -lprotoc -lre2 -lupb  -lz

win32: LIBS += -lws2_32

#here I will make use of pkg-config to get the list of dependencie of each libraries.
LIBS += `pkg-config --libs libplctag grpc++ protobuf`
