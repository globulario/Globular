QT -= gui

TEMPLATE = lib
CONFIG += staticlib

CONFIG += c++11

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
    globularserver.cpp

HEADERS += \
    globularserver.h \
    json.hpp

# Default rules for deployment.
unix {
    target.path = $$[QT_INSTALL_PLUGINS]/generic
}
!isEmpty(target.path): INSTALLS += target

win32:CONFIG(release, debug|release): LIBS += -L$$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/ -lGlobularRessourceClient
else:win32:CONFIG(debug, debug|release): LIBS += -L$$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/ -lGlobularRessourceClient
else:unix:!macx: LIBS += -L$$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/ -lGlobularRessourceClient

INCLUDEPATH += $$PWD/../../ressource/cpp/GlobularRessourceClient $$PWD/../../api/cpp/GlobularClient C:\Users\mm006819\grpc\include C:\Users\mm006819\grpc\third_party\protobuf\src C:\Users\mm006819\grpc\include

win32-g++:CONFIG(release, debug|release): PRE_TARGETDEPS += $$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/libGlobularRessourceClient.a
else:win32-g++:CONFIG(debug, debug|release): PRE_TARGETDEPS += $$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/libGlobularRessourceClient.a
else:win32:!win32-g++:CONFIG(release, debug|release): PRE_TARGETDEPS += $$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/GlobularRessourceClient.lib
else:win32:!win32-g++:CONFIG(debug, debug|release): PRE_TARGETDEPS += $$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/GlobularRessourceClient.lib
else:unix:!macx: PRE_TARGETDEPS += $$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/libGlobularRessourceClient.a
