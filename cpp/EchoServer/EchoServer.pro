TEMPLATE = app
CONFIG += console c++
CONFIG -= app_bundle
CONFIG -= qt

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

win32: LIBS += -lws2_32

win32: LIBS += -L$$PWD/../../../../../../../grpc/.build/third_party/protobuf/ -lprotobuf
win32: LIBS += -L$$PWD/../../../../../../../grpc/.build/ -lgrpc++ -lgpr -lgrpc -laddress_sorting -lgrpc++_reflection
#win32: LIBS += -L$$PWD/../../../../../../../grpc/.build/third_party/zlib/ -lzlibstatic
#win32: LIBS += -L$$PWD/../../../../../../../grpc/.build/third_party/cares/cares/lib -lcares
#win32: LIBS += -L$$PWD/../../../../../../../grpc/.build/third_party/abseil-cpp/absl/strings -labsl_cord -labsl_str_format_internal -labsl_strings -labsl_strings_internal
#win32: LIBS += -L$$PWD/../../../../../../../grpc/.build/third_party/abseil-cpp/absl/types -labsl_bad_any_cast_impl -labsl_bad_optional_access -labsl_bad_optional_access -labsl_bad_variant_access
#win32: LIBS += -L$$PWD/../../../../../../../grpc/.build/third_party/abseil-cpp/absl/base -labsl_base -labsl_dynamic_annotations -labsl_exponential_biased -labsl_log_severity -labsl_malloc_internal -labsl_periodic_sampler -labsl_raw_logging_internal -labsl_scoped_set_env -labsl_spinlock_wait -labsl_throw_delegate
#win32: LIBS += -L$$PWD/../../../../../../../grpc/.build/third_party/abseil-cpp/absl/synchronization -labsl_graphcycles_internal -labsl_synchronization
