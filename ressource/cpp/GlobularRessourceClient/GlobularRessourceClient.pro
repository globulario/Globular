QT -= gui

TEMPLATE = lib
CONFIG += staticlib

CONFIG += c++17

# The following define makes your compiler emit warnings if you use
# any Qt feature that has been marked deprecated (the exact warnings
# depend on your compiler). Please consult the documentation of the
# deprecated API in order to know how to port your code away from it.
DEFINES += QT_DEPRECATED_WARNINGS

# You can also make your code fail to compile if it uses deprecated APIs.
# In order to do so, uncomment the following line.
# You can also select to disable deprecated APIs only up to a certain version of Qt.
#DEFINES += QT_DISABLE_DEPRECATED_BEFORE=0x060000    # disables all the APIs deprecated before Qt 6.0.0

SOURCES += \
    globularressourceclient.cpp \
    ressource/ressource.grpc.pb.cc \
    ressource/ressource.pb.cc

HEADERS += \
    globularressourceclient.h \
    ressource/ressource.grpc.pb.h \
    ressource/ressource.pb.h

# Default rules for deployment.
unix {
    target.path = $$[QT_INSTALL_PLUGINS]/generic
}
!isEmpty(target.path): INSTALLS += target

INCLUDEPATH += $$PWD/../../../api/cpp/GlobularClient C:\Users\mm006819\grpc\include C:\Users\mm006819\grpc\third_party\protobuf\src
win32: LIBS += -L$$PWD/../../api/cpp/build-GlobularClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release -lGlobularClient
