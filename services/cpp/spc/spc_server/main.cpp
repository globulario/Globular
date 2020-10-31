#include <iostream>
#include <string>
#include <QtDebug>
#include <QString>
#include <QJsonDocument>
#include <QJsonArray>
#include <QJsonObject>

// The rpc service.
#include <grpcpp/grpcpp.h>
#include "spc/spcpb/spc.pb.h"
#include "spc/spcpb/spc.grpc.pb.h"

// Specific data structure.
#include "DonneesAnalyse.h"
#include "AnalyseurCSP.h"
#include "globularserver.h"
#include <iostream>
#include "cxxopts.hpp"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;
using namespace std;

class SpcServiceImpl final : public spc::SpcService::Service, public Globular::GlobularService
{
public:
    SpcServiceImpl(std::string id = "spc",
                   std::string name = "spc.SpcService",
                   std::string domain = "localhost",
                   std::string publisher_id = "localhost",
                   bool allow_all_origins = false,
                   std::string allowed_origins = "",
                   bool tls = false,
                   std::string version = "0.0.1",
                   unsigned int defaultPort = 10061, unsigned int defaultProxy = 10062):
         Globular::GlobularService(id, name, domain, publisher_id, allow_all_origins, allowed_origins, version, tls, defaultPort, defaultProxy ){
        // Set the proto path if is not already set.
        if(this->proto.length() == 0){
            ::spc::CreateAnalyseRqst request;
            this->proto = this->root + "/" + request.GetDescriptor()->file()->name();
            this->save();
        }

    }



    /**
    *   There is where the analsyse is done.
    */
    ::grpc::Status CreateAnalyse(::grpc::ServerContext* context, const ::spc::CreateAnalyseRqst* request, ::spc::CreateAnalyseRsp* response) override{
        // Here I will parse the json object from the input string.
        QJsonArray data = QJsonDocument::fromJson(QString::fromStdString(request->data()).toUtf8()).array();
        QJsonArray tests = QJsonDocument::fromJson(QString::fromStdString(request->tests()).toUtf8()).array();
        bool isPopulation = request->ispopulation();
        double tolzon = request->tolzon();
        double lotol = request->lotol();
        double uptol = request->uptol();
        QString tolType = QString::fromStdString(request->toltype()); // uncomment when it will corrected.

        /*if (parseError.error != QJsonParseError::NoError)
        {
            qDebug() << "Parse error: " << parseError.errorString();
            return grpc::Status(grpc::StatusCode::INTERNAL,
                           parseError.errorString().toStdString());
        }*/

        DonneesAnalyse analysedData;
        // Ici je vais convertir l'information reçu...
        for( QJsonArray::const_iterator it = data.begin(); it != data.end(); it++){
            // les lignes.
            QJsonArray row = (*it).toArray();
            QString serial = row.at(0).toString();
            double cote = row.at(1).toDouble();
            QString dateStr = row.at(2).toString();
            bool isActif = row.at(3).toBool();

            PieceInfo* info = new PieceInfo();
            info->serial = serial.toStdString();
            info->tolZone  = cote;
            info->creationDate = dateStr.toStdString();
            info->path = "";

            info->isSelect = isActif; // TODO change it to bool value...
            analysedData.addPieceInfo(info);
        }

        // Le type de tolerance.
        analysedData.setTolOption(tolType);

        // Set the test information.
        // K1
        analysedData.setState_test_K1(tests[0].toObject()["state"].toBool());
        analysedData.setTest_K1(tests[0].toObject()["value"].toDouble());
        // qDebug() <<"Test k1 value " << tests[0].toObject()["value"].toInt();

        // K2
        analysedData.setState_test_K2(tests[1].toObject()["state"].toBool());
        analysedData.setTest_K2(tests[1].toObject()["value"].toInt());
        // qDebug() <<"Test k2 value " << tests[1].toObject()["value"].toInt();

        // K3
        analysedData.setState_test_K3(tests[2].toObject()["state"].toBool());
        analysedData.setTest_K3(tests[2].toObject()["value"].toInt());
        // qDebug() <<"Test k3 value " << tests[2].toObject()["value"].toInt();

        // K4
        analysedData.setState_test_K4(tests[3].toObject()["state"].toBool());
        analysedData.setTest_K4(tests[3].toObject()["value"].toInt());
        // qDebug() <<"Test k4 value " << tests[3].toObject()["value"].toInt();

        // K5
        analysedData.setState_test_K5(tests[4].toObject()["state"].toBool());
        analysedData.setTest_K5(tests[4].toObject()["value"].toInt());
        //qDebug() <<"Test k5 value " << tests[4].toObject()["value"].toInt();

        // K6
        analysedData.setState_test_K6(tests[5].toObject()["state"].toBool());
        analysedData.setTest_K6(tests[5].toObject()["value"].toInt());
        //  qDebug() <<"Test k6 value " << tests[5].toObject()["value"].toInt();

        // K7
        analysedData.setState_test_K7(tests[6].toObject()["state"].toBool());
        analysedData.setTest_K7(tests[6].toObject()["value"].toInt());
        // qDebug() <<"Test k7 value " << tests[6].toObject()["value"].toInt();

        // K8
        analysedData.setState_test_K8(tests[7].toObject()["state"].toBool());
        analysedData.setTest_K8(tests[7].toObject()["value"].toInt());
        // qDebug() <<"Test k8 value " << tests[7].toObject()["value"].toInt();

        // Set if is an echantillon or a population.
        analysedData.setPopOuEchan(isPopulation);

        // Set tolerances
        analysedData.initToleranceInfo(tolzon, lotol, uptol);

        // Calculer l'analyseur
        AnalyseurCSP analyser;
        analyser.analyseComplete(analysedData);

        QJsonObject jsonObj;

        analysedData.write(jsonObj);

        // Here I will serialyse the result.
        QJsonDocument doc(jsonObj);
        QString strJson(doc.toJson(QJsonDocument::Compact));

        // Set the resonse message.
        response->set_result(strJson.toStdString());

        return ::grpc::Status::OK;
    }

    ::grpc::Status Stop(::grpc::ServerContext* /*context*/, const ::spc::StopRequest* /*request*/, ::spc::StopResponse* /*response*/) override {
        this->stop();
        return ::grpc::Status::OK;
    }

};

int main(int argc, char** argv)
{

    cxxopts::Options options("Statistic process control service", "A c++ gRpc service implementation");
    auto result = options.parse(argc, argv);

    // Instantiate a new server.
    SpcServiceImpl service;

    // Start the service.
    service.run(&service);

    return 0;
}
