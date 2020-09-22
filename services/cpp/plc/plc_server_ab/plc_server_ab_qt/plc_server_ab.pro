TEMPLATE = app
CONFIG += console c++11
CONFIG -= app_bundle
CONFIG -= qt

SOURCES += \
    ../../../plcpb/cpp/plc.grpc.pb.cc \
    ../../../plcpb/cpp/plc.pb.cc \
    ../PLC_server/PlcServiceImpl.cpp \
    ../PLC_server/main.cpp

HEADERS += \
    ../../../plcpb/cpp/plc.grpc.pb.h \
    ../../../plcpb/cpp/plc.pb.h \
    ../PLC_server/PlcServiceImpl.h \
    ../PLC_server/cxxopts.hpp \
    ../PLC_server/json.hpp

INCLUDEPATH += ../../../plcpb/cpp


LIBS += `pkg-config --libs libplctag grpc++ protobuf`
