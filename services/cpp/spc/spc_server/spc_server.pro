TEMPLATE = app
CONFIG += console c++17
CONFIG -= app_bundle
QT       += core

SOURCES += \
    ../../GlobularClient/globularclient.cpp \
    ../../GlobularServer/globularserver.cpp \
    ../../ressource/GlobularRessourceClient/globularressourceclient.cpp \
    ../../ressource/ressourcepb/ressource.grpc.pb.cc \
    ../../ressource/ressourcepb/ressource.pb.cc \
    ../spcpb/spc.grpc.pb.cc \
    ../spcpb/spc.pb.cc \
    AnalyseurCSP.cpp \
    DonneesAnalyse.cpp \
    Erreur.cpp \
    main.cpp \
    SousGroupe.cpp

HEADERS += \
    ../../GlobularClient/globularclient.h \
    ../../GlobularServer/globularserver.h \
    ../../ressource/GlobularRessourceClient/globularressourceclient.h \
    ../../ressource/ressourcepb/ressource.grpc.pb.h \
    ../../ressource/ressourcepb/ressource.pb.h \
    ../spcpb/spc.grpc.pb.h \
    ../spcpb/spc.pb.h \
    AnalyseurCSP.h \
    DonneesAnalyse.h \
    Erreur.h \
    SousGroupe.h


INCLUDEPATH += C:\Users\mm006819\boost_1_74_0 ../../
INCLUDEPATH +=  ../../ ../echopb ../../GlobularServer ../../GlobularClient ../../ressource/GlobularRessourceClient ../../ressource/ressourcepb

unix:!macx:INCLUDEPATH += /usr/local/include
win32:INCLUDEPATH += C:/msys64/mingw64/include

#here I will make use of pkg-config to get the list of dependencie of each libraries.
unix: LIBS += `pkg-config --libs grpc++ protobuf`

# Set the pkconfig.
win32: LIBS += -LC:/msys64/mingw64/lib -lgrpc++ -labsl_raw_hash_set -labsl_hashtablez_sampler -labsl_exponential_biased -labsl_hash -labsl_bad_variant_access -labsl_city -labsl_status -labsl_cord -labsl_bad_optional_access -labsl_str_format_internal -labsl_synchronization -labsl_graphcycles_internal -labsl_symbolize -labsl_demangle_internal -labsl_stacktrace -labsl_debugging_internal -labsl_malloc_internal -labsl_time -labsl_time_zone -labsl_civil_time -labsl_strings -labsl_strings_internal -labsl_throw_delegate -labsl_int128 -labsl_base -labsl_spinlock_wait -labsl_raw_logging_internal -labsl_log_severity -labsl_dynamic_annotations -lgrpc -laddress_sorting -lre2 -lupb -lcares -lz -labsl_raw_hash_set -labsl_hashtablez_sampler -labsl_exponential_biased -labsl_hash -labsl_bad_variant_access -labsl_city -labsl_status -labsl_cord -labsl_bad_optional_access -labsl_str_format_internal -labsl_synchronization -labsl_graphcycles_internal -labsl_symbolize -labsl_demangle_internal -labsl_stacktrace -labsl_debugging_internal -labsl_malloc_internal -labsl_time -labsl_time_zone -labsl_civil_time -labsl_strings -labsl_strings_internal -labsl_throw_delegate -labsl_int128 -labsl_base -labsl_spinlock_wait -labsl_raw_logging_internal -labsl_log_severity -labsl_dynamic_annotations -lgpr -labsl_str_format_internal -labsl_synchronization -labsl_graphcycles_internal -labsl_symbolize -labsl_demangle_internal -labsl_stacktrace -labsl_debugging_internal -labsl_malloc_internal -labsl_time -labsl_time_zone -labsl_civil_time -labsl_strings -labsl_strings_internal -labsl_throw_delegate -labsl_int128 -labsl_base -labsl_spinlock_wait -labsl_raw_logging_internal -labsl_log_severity -labsl_dynamic_annotations -lssl -lcrypto -lws2_32 -lgdi32 -lcrypt32  -limagehlp -lprotobuf -lgmp

