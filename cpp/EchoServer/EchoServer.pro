TEMPLATE = app
CONFIG += console c++11
CONFIG -= app_bundle
CONFIG -= qt

SOURCES += \
        echo/echopb/echo.grpc.pb.cc \
        echo/echopb/echo.pb.cc \
        echoserviceimpl.cpp \
        main.cpp

INCLUDEPATH += C:\Users\mm006819\grpc\include C:\Users\mm006819\grpc\third_party\protobuf\src

win32:CONFIG(release, debug|release): LIBS += -L$$PWD/../build-GlobularServer-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/ -lGlobularServer
else:win32:CONFIG(debug, debug|release): LIBS += -L$$PWD/../build-GlobularServer-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/ -lGlobularServer
else:unix:!macx: LIBS += -L$$PWD/../build-GlobularServer-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/ -lGlobularServer

win32:CONFIG(release, debug|release): LIBS += -L$$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/ -lGlobularRessourceClient
else:win32:CONFIG(debug, debug|release): LIBS += -L$$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/ -lGlobularRessourceClient
else:unix:!macx: LIBS += -L$$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/ -lGlobularRessourceClient


win32-g++:CONFIG(release, debug|release): PRE_TARGETDEPS += $$PWD/../build-GlobularServer-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/libGlobularServer.a
else:win32-g++:CONFIG(debug, debug|release): PRE_TARGETDEPS += $$PWD/../build-GlobularServer-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/libGlobularServer.a
else:win32:!win32-g++:CONFIG(release, debug|release): PRE_TARGETDEPS += $$PWD/../build-GlobularServer-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/GlobularServer.lib
else:win32:!win32-g++:CONFIG(debug, debug|release): PRE_TARGETDEPS += $$PWD/../build-GlobularServer-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/GlobularServer.lib
else:unix:!macx: PRE_TARGETDEPS += $$PWD/../build-GlobularServer-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/libGlobularServer.a

HEADERS += \
    echo/echopb/echo.grpc.pb.h \
    echo/echopb/echo.pb.h \
    echoserviceimpl.h

unix:!macx|win32: LIBS += -L$$PWD/../../../../../../../grpc/.build/ -lgrpc -lgrpc++

INCLUDEPATH += $$PWD/../../../../../../../grpc/.build
DEPENDPATH += $$PWD/../../../../../../../grpc/.build

win32:!win32-g++: PRE_TARGETDEPS += $$PWD/../../../../../../../grpc/.build/grpc.lib
else:unix:!macx|win32-g++: PRE_TARGETDEPS += $$PWD/../../../../../../../grpc/.build/libgrpc.a
unix:!macx|win32: LIBS += -L$$PWD/../../../../../../../grpc/.build/third_party/protobuf/ -lprotobuf

INCLUDEPATH += $$PWD/../../../../../../../grpc/.build/third_party/protobuf
DEPENDPATH += $$PWD/../../../../../../../grpc/.build/third_party/protobuf

win32:!win32-g++: PRE_TARGETDEPS += $$PWD/../../../../../../../grpc/.build/third_party/protobuf/protobuf.lib
else:unix:!macx|win32-g++: PRE_TARGETDEPS += $$PWD/../../../../../../../grpc/.build/third_party/protobuf/libprotobuf.a

INCLUDEPATH += $$PWD/../GlobularServer  $$PWD/../../ressource/cpp/GlobularRessourceClient $$PWD/../../api/cpp/GlobularClient


win32-g++:CONFIG(release, debug|release): PRE_TARGETDEPS += $$PWD/../build-GlobularServer-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/libGlobularServer.a
else:win32-g++:CONFIG(debug, debug|release): PRE_TARGETDEPS += $$PWD/../build-GlobularServer-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/libGlobularServer.a
else:win32:!win32-g++:CONFIG(release, debug|release): PRE_TARGETDEPS += $$PWD/../build-GlobularServer-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/GlobularServer.lib
else:win32:!win32-g++:CONFIG(debug, debug|release): PRE_TARGETDEPS += $$PWD/../build-GlobularServer-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/GlobularServer.lib
else:unix:!macx: PRE_TARGETDEPS += $$PWD/../build-GlobularServer-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/libGlobularServer.a

win32-g++:CONFIG(release, debug|release): PRE_TARGETDEPS += $$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/libGlobularRessourceClient.a
else:win32-g++:CONFIG(debug, debug|release): PRE_TARGETDEPS += $$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/libGlobularRessourceClient.a
else:win32:!win32-g++:CONFIG(release, debug|release): PRE_TARGETDEPS += $$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/GlobularRessourceClient.lib
else:win32:!win32-g++:CONFIG(debug, debug|release): PRE_TARGETDEPS += $$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/GlobularRessourceClient.lib
else:unix:!macx: PRE_TARGETDEPS += $$PWD/../../ressource/cpp/build-GlobularRessourceClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/libGlobularRessourceClient.a

win32:CONFIG(release, debug|release): LIBS += -L$$PWD/../../api/cpp/build-GlobularClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/ -lGlobularClient
else:win32:CONFIG(debug, debug|release): LIBS += -L$$PWD/../../api/cpp/build-GlobularClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/ -lGlobularClient
else:unix:!macx: LIBS += -L$$PWD/../../api/cpp/build-GlobularClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/ -lGlobularClient

win32-g++:CONFIG(release, debug|release): PRE_TARGETDEPS += $$PWD/../../api/cpp/build-GlobularClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/libGlobularClient.a
else:win32-g++:CONFIG(debug, debug|release): PRE_TARGETDEPS += $$PWD/../../api/cpp/build-GlobularClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/libGlobularClient.a
else:win32:!win32-g++:CONFIG(release, debug|release): PRE_TARGETDEPS += $$PWD/../../api/cpp/build-GlobularClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/release/GlobularClient.lib
else:win32:!win32-g++:CONFIG(debug, debug|release): PRE_TARGETDEPS += $$PWD/../../api/cpp/build-GlobularClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/debug/GlobularClient.lib
else:unix:!macx: PRE_TARGETDEPS += $$PWD/../../api/cpp/build-GlobularClient-Desktop_Qt_static_MinGW_w64_64bit_MSYS2/libGlobularClient.a

win32: LIBS += -lws2_32
