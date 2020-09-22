#include "DonneesAnalyse.h"
#include "AnalyseurCSP.h"
#include <QStringList>
#include <fstream>
#include <QDateTime>
#include <list>
#include <algorithm>
#include <QDebug>
#include <boost\math\distributions\weibull.hpp>

// Cette fontion permet de retrouver la précision
// d'un nombre a virgule... par exemple 123.45 revois 2.
int getPrecision(double val){
    std::stringstream ss;
    ss << val;

    std::string strVal = ss.str();
    size_t indexOfPoint = strVal.find(".");
    if(indexOfPoint != std::string::npos){
        return int(strVal.length() - indexOfPoint - 1);
    }

    return 0;
}

using std::list;

//Constructeur et Destructeur
DonneesAnalyse::DonneesAnalyse()
    : Variance(0)
{
    setTest_K1(3.00);
    setTest_K2(9);
    setTest_K3(6);
    setTest_K4(14);
    setTest_K5(2);
    setTest_K6(4);
    setTest_K7(15);
    setTest_K8(8);
    state_test_K1 = false;
    state_test_K2 = false;
    state_test_K3 = false;
    state_test_K4 = false;
    state_test_K5 = false;
    state_test_K6 = false;
    state_test_K7 = false;
    state_test_K8 = false;

    //Données des limites de contrôles des cartes I/MR
    LCS_X = 0.0f;
    LCI_X= 0.0f;
    LC_X= 0.0f;
    LCS_MR= 0.0f;
    LCI_MR= 0.0f;
    LC_MR= 0.0f;

    // Les valeurs fixes pour les limite de contrôle...
    LCX_fix=0.0f;
    LCS_X_fix=0.0f;
    LCI_X_fix=0.0f;
    date_LCX_fix=0.0f;

    moyene= 0.0f;
    ecart_type= 0.0f;
    etendue= 0.0f;
    cp= 0.0f;
    cpk= 0.0f;

    lotol= 0.0f;
    uptol= 0.0f;
    tolNominal= 0.0f;

    limiteInf= 0.0f;
    limiteSup= 0.0f;

    OpNo="0";

    isNormal = 0;
    Normality= 0.0f;
    Skewness= 0.0f;
    Kurtosis= 0.0f;
    CV= 0.0f;
    Variance= 0.0f;
    Median= 0.0f;
    AD= 0.0f;
    Pvalue= 0.0f;
    Z= 0.0f;
    shape = 0.0f;
    scale = 0.0f;

    //Si true, les calculs doivent etre faite pour une population
    //Si false, les calculs doivent etre faite pour un echantillion
    PopOuEchan= true;
    isStatic_ = true;

    // initialisation des date par défaut...
    QDate date = QDate::currentDate();
    this->DateFin = date.toString("yyyy-MM-dd").toStdString();
    this->DateDebut = "2006-01-01";

}

DonneesAnalyse::DonneesAnalyse(const DonneesAnalyse & other){
    setTest_K1(other.test_K1);
    setTest_K2(other.test_K2);
    setTest_K3(other.test_K3);
    setTest_K4(other.test_K4);
    setTest_K5(other.test_K5);
    setTest_K6(other.test_K6);
    setTest_K7(other.test_K7);
    setTest_K8(other.test_K8);
    setState_test_K1(other.state_test_K1);
    setState_test_K2(other.state_test_K2);
    setState_test_K3(other.state_test_K3);
    setState_test_K4(other.state_test_K4);
    setState_test_K5(other.state_test_K5);
    setState_test_K6(other.state_test_K6);
    setState_test_K7(other.state_test_K7);
    setState_test_K8(other.state_test_K8);
    setPath(other.path);

    setLCS_X(other.LCS_X);
    setLCI_X(other.LCI_X);
    setLC_X(other.LC_X);

    setLCI_X_Fix(other.LCI_X);
    setLCS_X_Fix(other.LCS_X);
    setLC_X_Fix(other.LC_X);
    this->setLC_X_Date(other.date_LCX_fix);

    setLCS_MR(other.LCS_MR);
    setLCI_MR(other.LCI_MR);
    setLC_MR(other.LC_MR);

    setMoyene(other.moyene);
    setEcartType(other.ecart_type);
    setEtendue(other.etendue);
    setCp(other.cp);
    setCpk(other.cpk);
    setLoTol(other.lotol);
    setUpTol(other.uptol);
    setTolNominal(other.tolNominal);
    setLimiteInf(other.limiteInf);
    setLimiteSup(other.limiteSup);
    setOpNo(QString::fromStdString(other.OpNo));
    setIsNormal(other.isNormal);
    setNormality(other.Normality);
    setSkewness(other.Skewness);
    setKurtosis(other.Kurtosis);
    setCV(other.CV);
    setVariance(other.Variance);
    setMedian(other.Median);
    setAD(other.AD);
    setPvalue(other.Pvalue);
    setZ(other.Z);
    setPopOuEchan(other.PopOuEchan);
    setIsStatic(other.isStatic_);
    DateFin = other.DateFin;
    DateDebut = other.DateDebut;

    NoTol = other.NoTol;
    NoFeat = other.NoFeat;
    NoModele = other.NoModele;
    tolOption = other.tolOption;
    Commentaire = other.Commentaire;

    // Copy des donnée d'item infos...
    this->piecesInfos.clear();
    for (int i = 0 ; i < other.piecesInfos.size() ; i++)
    {
        this->piecesInfos.push_back(new PieceInfo(*other.piecesInfos.at(i)));
    }

    this->unselect_workpaths.resize(other.unselect_workpaths.size());
    std::copy(other.unselect_workpaths.begin(), other.unselect_workpaths.end(),this->unselect_workpaths.begin());
}

DonneesAnalyse::~DonneesAnalyse()
{
    qDebug() << "delete analyse data";
    // Supression des information du vecteur...
    for(std::vector<PieceInfo*>::iterator it = this->piecesInfos.begin(); it != this->piecesInfos.end(); ++it){
        delete *it;
        *it = NULL;
    }
}

//Accesseurs des cartes IMR
double DonneesAnalyse::getLCS_MR()
{
    return LCS_MR;
}
double DonneesAnalyse::getLCI_MR()
{
    return LCI_MR;
}
double DonneesAnalyse::getLC_MR()
{
    return LC_MR;
}
double DonneesAnalyse::getLCS_X()
{
    return LCS_X;
}
double DonneesAnalyse::getLCI_X()
{
    return LCI_X;
}
double DonneesAnalyse::getLC_X()
{
    return LC_X;
}
double DonneesAnalyse::getLCS_X_Fix(){
    return this->LCS_X_fix;
}
double DonneesAnalyse::getLCI_X_Fix(){
    return LCI_X_fix;
}
double DonneesAnalyse::getLC_X_Fix(){
    return this->LCX_fix;
}
string DonneesAnalyse::getLC_X_Date(){
    return date_LCX_fix;
}
void DonneesAnalyse::setLCS_X_Fix(double LCS_X){
    LCS_X_fix = LCS_X;
}
void DonneesAnalyse::setLCI_X_Fix(double LCI_X){
    LCI_X_fix = LCI_X;
}
void DonneesAnalyse::setLC_X_Fix(double LC_X){
    this->LCX_fix = LC_X;
}
void DonneesAnalyse::setLC_X_Date(string date){
    this->date_LCX_fix = date;
}
double DonneesAnalyse::getMoyene()
{
    return moyene;
}
double DonneesAnalyse::getEtendue()
{
    return etendue;
}
double DonneesAnalyse::getEcartType()
{
    return ecart_type;
}
double DonneesAnalyse::getCp()
{
    return cp;
}
double DonneesAnalyse::getCpk()
{
    return cpk;
}
double DonneesAnalyse::getLoTol()
{
    return lotol;
}
double DonneesAnalyse::getUpTol()
{
    return uptol;
}
double DonneesAnalyse::getTolNominal()
{
    return tolNominal;
}
double DonneesAnalyse::getLimiteInf()
{
    return limiteInf;
}
double DonneesAnalyse::getLimiteSup()
{
    return limiteSup;
}
QString DonneesAnalyse::getNoTol()
{

    return QString::fromStdString(NoTol);
}
QString DonneesAnalyse::getNoFeat()
{
    return QString::fromStdString(NoFeat);
}
QString DonneesAnalyse::getOpNo()
{
    return QString::fromStdString(OpNo);
}
QString DonneesAnalyse::getNoModele()
{
    return QString::fromStdString(NoModele);
}
QString DonneesAnalyse::getDateDebut()
{
    return QString::fromStdString(DateDebut);
}
QString DonneesAnalyse::getDateFin()
{
    return QString::fromStdString(DateFin);
}
QString DonneesAnalyse:: getCommentaire()
{
    return QString::fromStdString(Commentaire);
}
int DonneesAnalyse::getIsNormal()
{
    return isNormal;
}
double DonneesAnalyse::getNormality()
{
    return Normality;
}
double DonneesAnalyse::getSkewness()
{
    return Skewness;
}
double DonneesAnalyse::getKurtosis()
{
    return Kurtosis;
}
double DonneesAnalyse::getCV()
{
    return CV;
}
double DonneesAnalyse::getVariance()
{
    return Variance;
}
double DonneesAnalyse::getMedian()
{
    return Median;
}
double DonneesAnalyse::getAD()
{
    return AD;
}
double DonneesAnalyse::getPvalue()
{
    return Pvalue;
}
double DonneesAnalyse::getZ()
{
    return Z;
}
bool DonneesAnalyse::getPopOuEchan()
{
    //Si true, les calculs doivent etre faite pour une population
    //Si false, les calculs doivent etre faite pour un echantillion
    return PopOuEchan;
}
double DonneesAnalyse::getScale(){
    return this->scale;
}
double DonneesAnalyse::getShape(){
    return this->shape;
}

//Accesseurs des paramètre de test pour l'analyse
double DonneesAnalyse::getTest_K1()
{
    return test_K1;
}
int DonneesAnalyse::getTest_K2()
{
    return test_K2;
}
int DonneesAnalyse::getTest_K3()
{
    return test_K3;
}
int DonneesAnalyse::getTest_K4()
{
    return test_K4;
}
int DonneesAnalyse::getTest_K5()
{
    return test_K5;
}
int DonneesAnalyse::getTest_K6()
{
    return test_K6;
}
int DonneesAnalyse::getTest_K7()
{
    return test_K7;
}
int DonneesAnalyse::getTest_K8()
{
    return test_K8;
}

//Acceusseurs de l'états des tests
bool DonneesAnalyse::getState_test_K1()
{
    return state_test_K1;
}
bool DonneesAnalyse::getState_test_K2()
{
    return state_test_K2;
}
bool DonneesAnalyse::getState_test_K3()
{
    return state_test_K3;
}
bool DonneesAnalyse::getState_test_K4()
{
    return state_test_K4;
}
bool DonneesAnalyse::getState_test_K5()
{
    return state_test_K5;
}
bool DonneesAnalyse::getState_test_K6()
{
    return state_test_K6;
}
bool DonneesAnalyse::getState_test_K7()
{
    return state_test_K7;
}
bool DonneesAnalyse::getState_test_K8()
{
    return state_test_K8;
}

//Mutateurs pour l'état des tests
void DonneesAnalyse::setState_test_K1(bool state_test_K1)
{
    this->state_test_K1 = state_test_K1;
}
void DonneesAnalyse::setState_test_K2(bool state_test_K2)
{
    this->state_test_K2 = state_test_K2;
}
void DonneesAnalyse::setState_test_K3(bool state_test_K3)
{
    this->state_test_K3 = state_test_K3;
}
void DonneesAnalyse::setState_test_K4(bool state_test_K4)
{
    this->state_test_K4 = state_test_K4;
}
void DonneesAnalyse::setState_test_K5(bool state_test_K5)
{
    this->state_test_K5 = state_test_K5;
}
void DonneesAnalyse::setState_test_K6(bool state_test_K6)
{
    this->state_test_K6 = state_test_K6;
}
void DonneesAnalyse::setState_test_K7(bool state_test_K7)
{
    this->state_test_K7 = state_test_K7;
}
void DonneesAnalyse::setState_test_K8(bool state_test_K8)
{
    this->state_test_K8 = state_test_K8;
}


//Mutateurs pour les paramètres de test pour l'analyse
void DonneesAnalyse::setTest_K1(double test_K1)
{
    this->test_K1 = test_K1;
}
void DonneesAnalyse::setTest_K2(int test_K2)
{
    this->test_K2 = test_K2;
}
void DonneesAnalyse::setTest_K3(int test_K3)
{
    this->test_K3 = test_K3;
}
void DonneesAnalyse::setTest_K4(int test_K4)
{
    this->test_K4 = test_K4;
}
void DonneesAnalyse::setTest_K5(int test_K5)
{
    this->test_K5 = test_K5;
}
void DonneesAnalyse::setTest_K6(int test_K6)
{
    this->test_K6 = test_K6;
}
void DonneesAnalyse::setTest_K7(int test_K7)
{
    this->test_K7 = test_K7;
}
void DonneesAnalyse::setTest_K8(int test_K8)
{
    this->test_K8 = test_K8;
}

//Mutateurs des cartes IMR
void DonneesAnalyse::setLCS_MR(double LCS_MR)
{
    this->LCS_MR = LCS_MR;
}
void DonneesAnalyse::setLCI_MR(double LCI_MR)
{
    this->LCI_MR = LCI_MR;
}
void DonneesAnalyse::setLC_MR(double LC_MR)
{
    this->LC_MR = LC_MR;
}
void DonneesAnalyse::setLCS_X(double LCS_X)
{
    this->LCS_X = LCS_X;
}
void DonneesAnalyse::setLC_X(double LC_X)
{
    this->LC_X = LC_X;
}
void DonneesAnalyse::setLCI_X(double LCI_X)
{
    this->LCI_X = LCI_X;
}
void DonneesAnalyse::setMoyene(double moyene)
{
    this->moyene = moyene;
}
void DonneesAnalyse::setEtendue(double etendue)
{
    this->etendue = etendue;
}
void DonneesAnalyse::setEcartType(double ecart_type)
{
    this->ecart_type = ecart_type;
}
void DonneesAnalyse::setCp(double cp)
{
    this->cp = cp;
}
void DonneesAnalyse::setCpk(double cpk)
{
    this->cpk = cpk;
}
void DonneesAnalyse::setLoTol(double lotol)
{
    this->lotol = lotol;
}
void DonneesAnalyse::setUpTol(double uptol)
{
    this->uptol = uptol;
}
void DonneesAnalyse::setTolNominal(double tolNominal)
{
    this->tolNominal = tolNominal;
}
void DonneesAnalyse::setLimiteInf(double limiteInf)
{
    this->limiteInf = limiteInf;
}
void DonneesAnalyse::setLimiteSup(double limiteSup)
{
    this->limiteSup = limiteSup;
}
void DonneesAnalyse::setNoTol(QString NoTol)
{
    this->NoTol = NoTol.toStdString();
}
void DonneesAnalyse::setNoFeat(QString NoFeat)
{
    this->NoFeat = NoFeat.toStdString();
}
void DonneesAnalyse::setOpNo(QString OpNo)
{
    this->OpNo = OpNo.toStdString();
}
void DonneesAnalyse::setNoModele(QString NoModele)
{
    this->NoModele = NoModele.toStdString();
}
void DonneesAnalyse::setDateDebut(QString DateDebut)
{
    this->DateDebut = DateDebut.toStdString();
}
void DonneesAnalyse::setDateFin(QString DateFin)
{
    this->DateFin = DateFin.toStdString();
}
void DonneesAnalyse::setCommentaire(QString Commentaire)
{
    this->Commentaire = Commentaire.toStdString();
}
void DonneesAnalyse::setIsNormal(int isNormal)
{
    this->isNormal = isNormal;
}
void DonneesAnalyse::setNormality(double Normality)
{
    this->Normality = Normality;
}
void DonneesAnalyse::setSkewness(double Skewness)
{
    this->Skewness = Skewness;
}
void DonneesAnalyse::setKurtosis(double Kurtosis)
{
    this->Kurtosis = Kurtosis;
}
void DonneesAnalyse::setCV(double CV)
{
    this->CV = CV;
}
void DonneesAnalyse::setVariance(double Variance)
{
    this->Variance = Variance;
}
void DonneesAnalyse::setMedian(double Median)
{
    this->Median = Median;
}
void DonneesAnalyse::setAD(double AD)
{
    this->AD = AD;
}
void DonneesAnalyse::setPvalue(double Pvalue)
{
    this->Pvalue = Pvalue;
}
void DonneesAnalyse::setZ(double Z)
{
    this->Z = Z;
}
void DonneesAnalyse::setPopOuEchan(bool PopOuEchan)
{
    //Si true, les calculs doivent etre faite pour une population
    //Si false, les calculs doivent etre faite pour un echantillion
    this->PopOuEchan = PopOuEchan;
}

vector<PieceInfo*>& DonneesAnalyse::getPieceInfosLst(){
    return this->piecesInfos;
}

//surcharge de l'opérateur =
DonneesAnalyse& DonneesAnalyse::operator =(const DonneesAnalyse &donnees)
{
    if(&donnees != this)
    {
        setTest_K1(donnees.test_K1);
        setTest_K2(donnees.test_K2);
        setTest_K3(donnees.test_K3);
        setTest_K4(donnees.test_K4);
        setTest_K5(donnees.test_K5);
        setTest_K6(donnees.test_K6);
        setTest_K7(donnees.test_K7);
        setTest_K8(donnees.test_K8);
        setState_test_K1(donnees.state_test_K1);
        setState_test_K2(donnees.state_test_K2);
        setState_test_K3(donnees.state_test_K3);
        setState_test_K4(donnees.state_test_K4);
        setState_test_K5(donnees.state_test_K5);
        setState_test_K6(donnees.state_test_K6);
        setState_test_K7(donnees.state_test_K7);
        setState_test_K8(donnees.state_test_K8);
        setPath(donnees.path);

        setLCS_X(donnees.LCS_X);
        setLCI_X(donnees.LCI_X);
        setLC_X(donnees.LC_X);
        setLCS_MR(donnees.LCS_MR);
        setLCI_MR(donnees.LCI_MR);
        setLC_MR(donnees.LC_MR);
        setLCI_X_Fix(donnees.LCI_X);
        setLCS_X_Fix(donnees.LCS_X);
        setLC_X_Fix(donnees.LC_X);
        setLC_X_Date(donnees.date_LCX_fix);
        setMoyene(donnees.moyene);
        setEcartType(donnees.ecart_type);
        setEtendue(donnees.etendue);
        setCp(donnees.cp);
        setCpk(donnees.cpk);
        setLoTol(donnees.lotol);
        setUpTol(donnees.uptol);
        setTolNominal(donnees.tolNominal);
        setLimiteInf(donnees.limiteInf);
        setLimiteSup(donnees.limiteSup);
        setOpNo(QString::fromStdString(donnees.OpNo));
        setIsNormal(donnees.isNormal);
        setNormality(donnees.Normality);
        setSkewness(donnees.Skewness);
        setKurtosis(donnees.Kurtosis);
        setCV(donnees.CV);
        setVariance(donnees.Variance);
        setMedian(donnees.Median);
        setAD(donnees.AD);
        setPvalue(donnees.Pvalue);
        setZ(donnees.Z);
        setPopOuEchan(donnees.PopOuEchan);
        setIsStatic(donnees.isStatic_);
        DateFin = donnees.DateFin;
        DateDebut = donnees.DateDebut;

        NoTol = donnees.NoTol;
        NoFeat = donnees.NoFeat;
        NoModele = donnees.NoModele;
        tolOption = donnees.tolOption;
        Commentaire = donnees.Commentaire;

        // Copy des donnée d'item infos...
        this->piecesInfos.clear();
        for (int i = 0 ; i < donnees.piecesInfos.size() ; i++)
        {
            this->piecesInfos.push_back(new PieceInfo(*donnees.piecesInfos.at(i)));
        }

        this->unselect_workpaths.resize(donnees.unselect_workpaths.size());
        std::copy(donnees.unselect_workpaths.begin(), donnees.unselect_workpaths.end(),this->unselect_workpaths.begin());
    }
    return *this;
}

string DonneesAnalyse::getPath(){
    return this->path;
}

void DonneesAnalyse::setPath(string path){
    this->path = path;
}

void DonneesAnalyse::clearInfos(){
    this->piecesInfos.clear();
}

void DonneesAnalyse::addPieceInfo(PieceInfo* info){
    // Ajoute simplement une piece a la liste.
    this->piecesInfos.push_back(info);
}

// La fonction qui permet de mettre en ordre le vecteur des item info...
bool isNewer(PieceInfo* d1, PieceInfo* d2){
    QDate date1 = QDate::fromString(QString::fromStdString(d1->creationDate), "yyyy-MM-dd");
    QDate date2 = QDate::fromString(QString::fromStdString(d2->creationDate), "yyyy-MM-dd");
    return date1 > date2;
}

QString DonneesAnalyse::getTolOption(){
    return QString::fromStdString(this->tolOption);
}

void DonneesAnalyse::setTolOption(QString opTol){
    this->tolOption = opTol.toStdString();
}

struct isEqual
{
    // members.
    std::string id;

    // Constructor
    isEqual(const std::string& s): id(s){

    }

    // Return true if it has the same id.
    bool operator()(PieceInfo* p)
    {
        return p->serial == id;
    }
};

/**
 * @brief DonneesAnalyse::initItemInfo
 * Initialyse item infos from a tow dimensional array of data.
 */
void DonneesAnalyse::initItemInfo(){

    // Calcul de l'evolution des indices...
    calculEvolution(this->piecesInfos.size());
}

inline int getSignificativePosition(double val){
    std::stringstream ss;
    ss << val;
    std::string str = ss.str();;
    std::size_t found = str.find(".");
    if(found!=std::string::npos){
        return int(str.length() - found);
    }

    return 0;
}

/**
 * Cette fonction calcul la sommation des donnée majoré par leur précision.
 * De cette manière il n'y a pas d'imprécision liée au virgule...
 *
 * @brief DonneesAnalyse::getDonneesIndividuelleSummation
 * @return
 */
int DonneesAnalyse::getDonneesIndividuelleSummation(){
    int precision = getDonneesIndividuellePrecision();
    int factor = pow(10, precision);
    int sum = 0;

    for(vector<double>::const_iterator it = getDonneesIndividuelle().begin(); it != getDonneesIndividuelle().end(); ++it){
        double f = *it;
        f *= factor;

        // Round to the upped value if the leading digit is >= .5
        if(f > 0){
            f += 0.5;
        }else if(f < 0){
            f -= 0.5;
        }
        sum += (int) f;
    }

    return sum;
}

double DonneesAnalyse::getDonneesIndividuellePrecision(){
    int precision = 0;
    for(vector<double>::const_iterator it = getDonneesIndividuelle().begin(); it != getDonneesIndividuelle().end(); ++it){
        int tmp = getPrecision(*it);
        if(tmp > precision){
            precision = tmp;
        }
    }
    return precision;
}

const vector<double>&  DonneesAnalyse::getDonneesIndividuelle(){

    // Dans le cas ou les valeur existe...
    if(!donneesIndividuelle.empty()){
        return donneesIndividuelle;
    }
    donneesIndividuelle.clear();
    for(std::vector<PieceInfo*>::iterator it2 = piecesInfos.begin();it2 != piecesInfos.end(); ++it2){
        if((*it2)->isSelect){
            donneesIndividuelle.push_back((*it2)->tolZone);
        }
    }

    // Calcul des parametre...
    if(!this->isSupposedToBeNormal()){
        calculWeibullParameter(shape, scale, donneesIndividuelle, tolNominal);
    }
    return donneesIndividuelle;
}

string DonneesAnalyse::getNote(){
    return this->note;
}

void DonneesAnalyse::setNote(string& note){
    this->note = note;
}

QString DonneesAnalyse::getDescription(){
    QString description;
    description = QString::fromStdString(this->tolOption) + "\n";
    description += QString::fromStdString(this->NoFeat) + ":";
    description += QString::fromStdString(this->NoTol) + "\n";
    description +=  QString::number(this->tolNominal) +  " ";
    // Dans ce cas je peux
    if(abs(this->uptol) == abs(this->lotol)){
        description +=  QChar(0x00B1) + QString::number(abs(this->uptol));
    }
    else{
        description += QString::number(this->lotol) + " +" + QString::number(this->uptol);
    }
    return description;
}

void  DonneesAnalyse::initToleranceInfo(double tolzon, double lotol, double uptol)
{
    // Set tolerances.
    setLoTol(lotol);
    setUpTol(uptol);
    setTolNominal(tolzon);
    // Set limits.
    setLimiteInf(getTolNominal() + getLoTol());
    setLimiteSup(getTolNominal() + getUpTol());

}

void DonneesAnalyse::calculEvolution(unsigned int size = 30){

    // Reinitialisation ds valeur...
    evolution_moyennes.clear();
    evolution_cp.clear();
    evolution_cpk.clear();

    // Un tableau temporaire de pièce...
    list<PieceInfo*> echantillons;

    // Je dois donc réévaluer toutes les moyennes...
    for(vector<PieceInfo*>::reverse_iterator iter = piecesInfos.rbegin(); iter != piecesInfos.rend(); ++iter){

        // je sauvegarde la piece dans le tableau...
        echantillons.push_back(*iter);

        // Je supprime la plus ancienne piece...
        if(echantillons.size() > size){
            //echantillons.front();
            echantillons.pop_front();
        }

        // Ici si le stack contient 30 pièce...
        if(echantillons.size() == size){
            // Dans ce cas si je calcul les indices...
            DonneesAnalyse donneesPourAnalyse;
            donneesIndividuelle.clear();

            for(list<PieceInfo*>::const_iterator it = echantillons.begin();it != echantillons.end(); ++it){
                donneesPourAnalyse.donneesIndividuelle.push_back((*it)->tolZone);
            }

            donneesPourAnalyse.setUpTol(this->getUpTol());
            donneesPourAnalyse.setLoTol(this->getLoTol());
            donneesPourAnalyse.setTolNominal(this->getTolNominal());
            donneesPourAnalyse.setLimiteInf(this->getLimiteInf());
            donneesPourAnalyse.setLimiteSup(this->getLimiteSup());

            double _shape = 0;
            double _scale = 0;
            calculWeibullParameter(_shape, _scale, donneesPourAnalyse.donneesIndividuelle, donneesPourAnalyse.tolNominal);
            donneesPourAnalyse.setScale(_scale);
            donneesPourAnalyse.setShape(_shape);


            // Calcul de la moyenne et cp/cpk...
            AnalyseurCSP analyseur;
            analyseur.calculMoyenne(donneesPourAnalyse);
            analyseur.calculVariance(donneesPourAnalyse);
            analyseur.calculEcartType(donneesPourAnalyse);
            analyseur.calculCpCpk(donneesPourAnalyse);

            // Ajout des indice...
            evolution_moyennes.push_back(donneesPourAnalyse.getMoyene());
            evolution_cp.push_back(donneesPourAnalyse.getCp());
            evolution_cpk.push_back(donneesPourAnalyse.getCpk());
        }
    }
}

// Les tendance...
Tendance  DonneesAnalyse::getTendanceMoyenne(){
    // La tendence de la moyenne ce calcul en fonction des deux dernière moyennes...
    if(this->evolution_moyennes.size() >= 2){
        double lastResult = evolution_moyennes[evolution_moyennes.size() - 1];
        double previousResult = evolution_moyennes[evolution_moyennes.size() - 2];
        double val1 = abs(this->tolNominal - lastResult);
        double val2 = abs(this->tolNominal - previousResult);
        if(val1 < val2){
            return Convergent;
        }
        if(val1 > val2){
            return Divergent;
        }
    }

    // Ici la valeur est égale.
    return Neutre;
}

Tendance DonneesAnalyse::getTendanceCp(){
    // La tendence de la moyenne ce calcul en fonction des deux dernière moyennes...
    if(this->evolution_cp.size() >= 2){
        double lastResult = evolution_cp[evolution_cp.size() - 1];
        double previousResult = evolution_cp[evolution_cp.size() - 2];
        double variation = (1 - (previousResult / lastResult))*100.0f;
        if(abs(variation) > 5){
            if(variation > 0){
                return Convergent;
            }
            if(variation < 0){
                return Divergent;
            }
        }
    }

    // Ici la valeur est égale.
    return Neutre;
}

Tendance DonneesAnalyse::getTendanceCpk(){
    // La tendence de la moyenne ce calcul en fonction des deux dernière moyennes...
    if(this->evolution_cpk.size() >= 2){
        double lastResult = evolution_cpk[evolution_cpk.size() - 1];
        double previousResult = evolution_cpk[evolution_cpk.size() - 2];
        double variation = (1 - (previousResult / lastResult))*100.0f;
        if(abs(variation) > 5){
            if(variation > 0){
                return Convergent;
            }
            if(variation < 0){
                return Divergent;
            }
        }
    }

    // Ici la valeur est égale.
    return Neutre;
}

bool DonneesAnalyse::isTendanceCpkAcceptable(){
    if (evolution_cpk.size() > 0)
    {
        if(cpk > 1.33){
            return true;
        }
    }
    return false;
}

bool DonneesAnalyse::isTendanceCpAcceptable(){
    if (evolution_cp.size() > 0)
    {
        if(cp > 2.0){
            return true;
        }
    }
    return false;
}

QVector<int> DonneesAnalyse::getFailedTests()
{
    return failedTests;
}

void DonneesAnalyse::setFailedTests(QVector<int> value)
{
    failedTests = value;
}

bool DonneesAnalyse::isStatic(){
    return this->isStatic_;
}

void DonneesAnalyse::setIsStatic(bool val){
    this->isStatic_ = val;
}

void DonneesAnalyse::update(){
    // Intialisation des donnée...
    if(!isStatic() || this->piecesInfos.empty()){
        initItemInfo();
        AnalyseurCSP analyseur;
        analyseur.calculIMR(*this);
    }
}

void DonneesAnalyse::update(std::string serial, bool isSelect){
    // Dans ce cas si je vais retrouver la liste
    // Mise a jour pour le prochain access...
    this->donneesIndividuelle.clear();

    // Mise a jour de l'element dans le vecteur d'item infos.
    std::vector<PieceInfo*>::iterator it = std::find_if(piecesInfos.begin(), piecesInfos.end(), isEqual(serial));
    if(it != piecesInfos.end()){
        (*it)->isSelect = isSelect;

        QString NoTol = this->getNoTol();
        QString NoFeat= this->getNoFeat();
        QString OpNo= this->getOpNo().split(";").at(1).split(":").at(0).trimmed();
        QString NoModele= this->getNoModele();
        QString tolOption= this->getTolOption();

        QStringList splits = QString::fromStdString((*it)->path).replace("\\\\MON-FILER-01\\Data\\Cmm\\1GeoPro\\", "").split("\\");
        QString subdir = "\\" + splits.at(splits.size()-2);

        // Ici je vais mettre a jour les pieces de l'échantillon...
        //qDebug() << "Item " + QString::fromStdString(serial) << " is selected!";
        // qDebug() << " " << NoTol << " " << NoFeat << " " << OpNo << " " <<  NoModele << " " << tolOption << " " << subdir;

        int val = 0;
        if(!isSelect){
            val = 1;
        }
    }
}

// Les fonction suivante servent a resoudre l'equation
double getFactor(double B, const std::vector<double>& data){
    double sumXiB = 0.0f; double sumLnXi = 0.0f; double sumXiBLnXi = 0.0f;
    int size = 0;
    for(std::vector<double>::const_iterator it = data.begin(); it != data.end(); ++it){
        double val = *it;
        if (val != 0)
        {
            double xiB = pow(val, B);
            double lnXi = log(val);
            double xiBLnXi = xiB * lnXi;
            sumXiB += xiB;
            sumLnXi += lnXi;
            sumXiBLnXi += xiBLnXi;
            ++size;
        }
    }

    double result = (sumXiBLnXi/sumXiB) - (1.0f/B) - (1.0f/size)*sumLnXi;
    return result;
}

// La valeur de départ est toujours la valeur la plus petite...
void getFactor(double increment, double& startValue, const std::vector<double>& data){

    double value = getFactor(startValue, data);

    while(value < 0){
        startValue += 1/increment;
        value = getFactor(startValue, data);
    }

    // La precision peut etre augmenter au besoin.
    if(increment < 10000){
        startValue -= (1/increment);
        getFactor(increment * 10, startValue, data);
    }
}

void DonneesAnalyse::calculWeibullParameter(double& shape, double& scale, std::vector<double> data, const double& tolNom){

    if(data.empty()){
        return;
    }

    // Je ramene toute les valeur en terme de deviance...
    int zeroValue = 0;
    for(std::vector<double>::iterator it0 = data.begin(); it0 != data.end(); ++it0){
        *it0 = abs(*it0 - tolNom);
        if (*it0 == 0)
        {
            zeroValue++;
        }
    }

    // True = incremente...
    getFactor(1.0f, shape, data);

    // Maintenant que je connais le facteur de forme je vais calculer le facteur de scale.
    double sumXiB = 0.0f;
    int size = 0;
    for(std::vector<double>::const_iterator it = data.begin(); it != data.end(); ++it){
        double val = *it;
        if (val != 0)
        {
            double xiB = pow(val, shape);
            sumXiB += xiB;
            ++size;
        }
    }

    if (shape <= 0)
    {
        shape = 1;
    }

    if (size == 0)
    {
        scale = 1;
    }
    else
    {
        scale = pow((sumXiB/size), (1/shape));
    }
}



void DonneesAnalyse::clearPieces()
{
    this->donneesIndividuelle.clear();
}


//retourne si la courbe devrait être normal ou non
bool DonneesAnalyse::isSupposedToBeNormal()
{
    //si la tolérence inférieur est a 0 et la tolérence supérieur est différente c'est supposé être une weibull...
    if (this->getLoTol() == 0 && this->getUpTol() > 0)
        return false;
    else
        return true;
}

//set shape
void DonneesAnalyse::setShape(double _shape)
{
    this->shape = _shape;
}
void DonneesAnalyse::setScale(double _scale)
{
    this->scale = _scale;
}

void DonneesAnalyse::addUnselectWorkPath(QString val){
    std::vector<std::string>::iterator findIter = std::find(this->unselect_workpaths.begin(), this->unselect_workpaths.end(), val.toStdString());
    if(findIter != this->unselect_workpaths.end()){
        this->unselect_workpaths.push_back(val.toStdString());
    }
}

void DonneesAnalyse::removeUnselectWorkPaht(QString val){
    std::vector<std::string>::iterator findIter = std::find(this->unselect_workpaths.begin(), this->unselect_workpaths.end(), val.toStdString());
    if(findIter != this->unselect_workpaths.end()){
        this->unselect_workpaths.erase(findIter);
    }
}

void DonneesAnalyse::setUnselectWorkPath(QList<QString>& lst){
    this->unselect_workpaths.clear();
    for(int i=0; i < lst.size(); ++i){
        this->unselect_workpaths.push_back(lst.at(i).toStdString());
    }
}

QList<QString>* DonneesAnalyse::getUnselectWorkPath()
{
    QList<QString>* lst = new QList<QString>();
    std::vector<std::string>::iterator iter = unselect_workpaths.begin();
    while (iter != this->unselect_workpaths.end()){
        lst->append(QString::fromStdString(*iter));
        iter++;
    }

    return lst;
}

// unmarshal
void DonneesAnalyse::read(const QJsonObject &json)
{
    // Set test values.
    this->test_K1 = json["test_K1"].toDouble();
    this->test_K2 = json["test_K2"].toInt();
    this->test_K3 = json["test_K3"].toInt();
    this->test_K4 = json["test_K4"].toInt();
    this->test_K5 = json["test_K5"].toInt();
    this->test_K6 = json["test_K6"].toInt();
    this->test_K7 = json["test_K7"].toInt();
    this->test_K8 = json["test_K8"].toInt();

    this->state_test_K1 = json["test_K1"].toBool();
    this->state_test_K2 = json["test_K2"].toBool();
    this->state_test_K3 = json["test_K3"].toBool();
    this->state_test_K4 = json["test_K4"].toBool();
    this->state_test_K5 = json["test_K5"].toBool();
    this->state_test_K6 = json["test_K6"].toBool();
    this->state_test_K7 = json["test_K7"].toBool();
    this->state_test_K8 = json["test_K8"].toBool();

    this->LCS_X = json["LCS_X"].toDouble();
    this->LCI_X = json["LCI_X"].toDouble();
    this->LC_X = json["LC_X"].toDouble();

    this->LCX_fix = json["LCX_fix"].toDouble();
    this->LCS_X_fix = json["LCS_X_fix"].toDouble();
    this->LCI_X_fix = json["LCI_X_fix"].toDouble();
    this->date_LCX_fix = json["date_LCX_fix"].toString().toStdString();

    this->LCS_MR = json["LCS_MR"].toDouble();
    this->LCI_MR =  json["LCI_MR"].toDouble();
    this->LC_MR = json["LC_MR"].toDouble();

    this->tolOption = json["tolOption"].toString().toStdString();
    this->moyene = json["moyene"].toDouble();
    this->ecart_type = json["ecart_type"].toDouble();
    this->etendue = json["etendue"].toDouble();
    this->cp = json["cp"].toDouble();
    this->cpk = json["cpk"].toDouble();
    this->isNormal = json["isNormal"].toBool();
    this->Normality = json["Normality"].toDouble();
    this->Skewness = json["Skewness"].toDouble();
    this->Kurtosis = json["Kurtosis"].toDouble();
    this->CV =  json["CV"].toDouble();
    this->Variance = json["Variance"].toDouble();
    this->Median = json["Median"].toDouble();
    this->AD = json["AD"].toDouble();
    this->Pvalue = json["Pvalue"].toDouble();
    this->Z =  json["Z"].toDouble();
    this->PopOuEchan = json["PopOuEchan"].toBool();

    this->lotol = json["lotol"].toDouble();
    this->uptol = json["uptol"].toDouble();
    this->tolNominal = json["tolNominal"].toDouble();

    this->limiteInf = json["limiteInf"].toDouble();
    this->limiteSup = json["limiteSup"].toDouble();

    this->DateDebut = json["DateDebut"].toString().toStdString();
    this->DateFin = json["DateFin"].toString().toStdString();
    this->Commentaire = json["Commentaire"].toString().toStdString();

    // Specific
    this->OpNo =  json["OpNo"].toString().toStdString();
    this->NoModele = json["NoModele"].toString().toStdString();
    this->NoTol = json["NoTol"].toString().toStdString();
    this->NoFeat = json["NoFeat"].toString().toStdString();

    // Distribution shape and scale.
    this->shape = json["shape"].toDouble();
}

// Utilisé pour calculer les valeur de l'histogram.
const double COURBE_NORMAL_SCALE_STEP = 0.30;

//méthode qui retourne un facteur pour mettre les histogrammes au niveaux
double scaleCourbeNormale(int scaleHistoCurrentDepth, const QVector<double> xInterval, const QVector<double> yIntervalOriginal, const boost::math::normal_distribution<> dist, const double moyenneParent, const double scaleParent)
{
    scaleHistoCurrentDepth++;
    if (scaleHistoCurrentDepth > 50)
    {
        scaleHistoCurrentDepth = 0;
        return 1;
    }

    QVector<double> pointsX(xInterval);
    QVector<double> pointsY(yIntervalOriginal);

    double sommeDistance = 0.0;
    double sommeHauteHisto = 0.0;

    //on applique un scale plus grand que celui du parent
    for (int i = 0 ; i < pointsY.size(); i++)
    {
        pointsY[i] *= (scaleParent + COURBE_NORMAL_SCALE_STEP);
    }

    //on calcul la distance entre la courbe normal et les histogrammes
    for (int i = 0 ; i < pointsY.size() ; i++)
    {
        sommeDistance += (boost::math::pdf(dist,pointsX[i]) - pointsY[i]);
        sommeHauteHisto += pointsY[i];
    }
    double moyenneDistance = sommeDistance / pointsY.size();


    //si la moyenne s'éloigne de la moyenne du parent on arrête la récursion
    if (moyenneParent < 0 && moyenneDistance < moyenneParent)
    {
        scaleHistoCurrentDepth = 0;
        return 1;
    }
    else if (moyenneParent > 0 && moyenneDistance > moyenneParent)
    {
        scaleHistoCurrentDepth = 0;
        return 1;
    }


    if (moyenneDistance > 1)
    {
        //TODO AJOUTER 0.05% AU Y dES POINTS ET RAPPELER SCALECOURBENORMAL
        return COURBE_NORMAL_SCALE_STEP + scaleCourbeNormale(scaleHistoCurrentDepth, pointsX,yIntervalOriginal,dist,moyenneDistance,scaleParent + COURBE_NORMAL_SCALE_STEP);
    }
    else if (moyenneDistance < -1)
    {
        //TODO AJOUTER -0.05% AU Y dES POINTS ET RAPPELER SCALECOURBENORMAL
        return -1 * COURBE_NORMAL_SCALE_STEP + scaleCourbeNormale(scaleHistoCurrentDepth, pointsX,yIntervalOriginal,dist,moyenneDistance,scaleParent - COURBE_NORMAL_SCALE_STEP);
    }
    else
    {
        scaleHistoCurrentDepth = 0;
        return 1;
    }
}

// marshal
void DonneesAnalyse::write(QJsonObject &json) const
{
    // Test results.
    json["test_K1"] = this->test_K1;
    json["test_K2"] = this->test_K2;
    json["test_K3"] = this->test_K3;
    json["test_K4"] = this->test_K4;
    json["test_K5"] = this->test_K5;
    json["test_K6"] = this->test_K6;
    json["test_K7"] = this->test_K7;
    json["test_K8"] = this->test_K8;

    // Test results states
    json["test_K1"] = this->state_test_K1;
    json["test_K2"] = this->state_test_K2;
    json["test_K3"] = this->state_test_K3;
    json["test_K4"] = this->state_test_K4;
    json["test_K5"] = this->state_test_K5;
    json["test_K6"] = this->state_test_K6;
    json["test_K7"] = this->state_test_K7;
    json["test_K8"] = this->state_test_K8;

    // Control limits.
    json["LCS_X"] = this->LCS_X;
    json["LCI_X"] = this->LCI_X;
    json["LC_X"] = this->LC_X;

    json["LCX_fix"] = this->LCX_fix;
    json["LCS_X_fix"] = this->LCS_X_fix;
    json["LCI_X_fix"] = this->LCI_X_fix;
    json["date_LCX_fix"] = QString::fromStdString(this->date_LCX_fix);

    json["LCS_MR"] = this->LCS_MR;
    json["LCI_MR"] = this->LCI_MR;
    json["LC_MR"] = this->LC_MR;

    json["tolOption"] = QString::fromStdString(this->tolOption);
    json["moyene"] = this->moyene;
    json["ecart_type"] = this->ecart_type;
    json["etendue"] = this->etendue;
    json["cp"] = this->cp;
    json["cpk"] = this->cpk;
    json["isNormal"] = this->isNormal;
    json["Normality"] = this->Normality;
    json["Skewness"] = this->Skewness;
    json["Kurtosis"] = this->Kurtosis;
    json["CV"] = this->CV;
    json["Variance"] = this->Variance;
    json["Median"] = this->Median;
    json["AD"] = this->AD;
    json["Pvalue"] = this->Pvalue;
    json["Z"] = this->Z;
    json["PopOuEchan"] = this->PopOuEchan;

    json["lotol"] = this->lotol;
    json["uptol"] = this->uptol;
    json["tolNominal"] = this->tolNominal;

    json["limiteInf"] = this->limiteInf;
    json["limiteSup"] = this->limiteSup;

    json["DateDebut"] = QString::fromStdString(this->DateDebut);
    json["DateFin"] = QString::fromStdString(this->DateFin);
    json["Commentaire"] = QString::fromStdString(this->Commentaire);

    QJsonArray state_X;
    for(size_t i=0; i <  this->state_X.size(); i++){
        QJsonObject obj;
        this->state_X[i].write(obj);
        state_X.append(obj);
    }
    json["state_X"] = state_X;

    QJsonArray state_MR;
    for(size_t i=0; i <  this->state_MR.size(); i++){
        QJsonObject obj;
        this->state_MR[i].write(obj);
        state_MR.append(obj);
    }
    json["state_MR"] = state_MR;

    QJsonArray coteZDonneeIndivitual;
    for(size_t i=0; i < this->coteZDonneeIndivitual.size(); i++){
        coteZDonneeIndivitual.append(this->coteZDonneeIndivitual[i]);
    }
    json["coteZDonneeIndivitual"] = coteZDonneeIndivitual;


    // Sous group imr...
    QJsonArray sousGroupe_IMR;
    for(size_t i=0; i < this->sousGroupe_IMR.size(); i++){
        QJsonObject obj;
        this->sousGroupe_IMR[i].write(obj);
        sousGroupe_IMR.append(obj);
    }
    json["sousGroupe_IMR"] = sousGroupe_IMR;

    // specific.
    json["NoTol"] = QString::fromStdString(this->NoTol);
    json["NoFeat"] = QString::fromStdString(this->NoFeat);
    json["OpNo"] = QString::fromStdString(this->OpNo);
    json["NoModele"] = QString::fromStdString(this->NoModele);

    json["shape"] = this->shape;
    json["scale"] = this->scale;

    double m = this->moyene;
    double s = this->ecart_type;
    double n = 0;

    // Copie of local values.
    QJsonArray items;

    std::vector<double> donneesIndivitual;
    for(std::vector<PieceInfo*>::const_iterator it = this->piecesInfos.cbegin(); it != this->piecesInfos.cend(); it++){
        if((*it)->isSelect == true){
            n++;
            donneesIndivitual.push_back((*it)->tolZone);
            QJsonObject item;
            item["serial"] = QString::fromStdString((*it)->serial);
            item["value"] = (*it)->tolZone;
            item["date"] = QString::fromStdString((*it)->creationDate);
            items.push_back(item);
        }
    }

    // Return information about input items.
    json["items"] = items;

    // met les donnees en ordre.
    std::sort(donneesIndivitual.begin(), donneesIndivitual.end());

    int nbi = 1 + (10*log(n))/3;
    double maxY = 0;
    double minY = 0;

    // Les donnés de la courbe.
    QVector<double> yPoints;
    QVector<double> xPoints;

    // Les donnée de l'histogram.
    QVector<double> pointsY;
    QVector<double> pointsX;

    if(!(this->lotol == 0 && this->uptol> 0)){
        // Normal distribution.

        // Point de la courbe normal.
        double range = (m+4*s)-(m-4*s);
        double nn = 1000;
        double ecart = range/nn;

        xPoints.push_back(m-4*s);
        //Calcul les points x
        for(int i = 0; i < nn; i++)
        {
            xPoints.push_back(xPoints[i]+ecart);
        }

        //Sort
        std::sort(xPoints.begin(), xPoints.end());

        //création de la distribution normal
        if(s == 0){
            // s doit etre plus grand que zero!!!
            s = 1;
        }

        boost::math::normal_distribution<> myNormal(m,s);

        //Calcul les donnees de y
        for(double i = m-(range/2); i < m+(range/2); i += range/nn)
        {
            yPoints.push_back(boost::math::pdf(myNormal,i));
        }

        minY=yPoints[0];
        for(int i = 0; i < nn; i++)
        {
            //Conserve la plus grande coordonnee Y
            if(yPoints[i] > maxY)
            {
                maxY = yPoints[i];
            }

            //Conserve la plus petite coordonnee Y
            if(yPoints[i] < minY)
            {
                minY = yPoints[i];
            }
        }
        // Normal histogram values.
        double maxAutreY = 0;
        //grandeur des division
        double ecartItv = this->etendue/nbi;

        vector<double> intervalle;
        for(int i = 0; i < nbi+1; i++)
        {
            intervalle.push_back(donneesIndivitual[0]+i*ecartItv);
        }

        if(ecartItv != 0 && intervalle.size() != 0)
        {
            //hauteur des y
            for(int i = 0; i < nbi; i++)
            {
                //mets tous a zero
                pointsY.push_back(0);
            }

            //calcul des fréquences pour chaque intervalle
            for(int i = 0; i < n; i++)
            {
                int x = ((donneesIndivitual[i] - donneesIndivitual[0]) / ecartItv) - 0.5;
                pointsY[x]++;

                if(pointsY[x] > maxAutreY)
                {
                    maxAutreY = pointsY[x];
                }
            }

            //on set les X des histogrammes
            for(int i = 0; i < nbi && i < intervalle.size(); i++)
            {
                pointsX.push_back(intervalle[i]+ecartItv/2);
            }

        }

    }else if(this->scale > 0 && this->shape > 0){

        double range = (m+4*s)-(m-4*s);
        double nn = 1000;
        double ecart = range/nn;

        xPoints.push_back(m-4*s);

        //Calcul les points x
        for(int i = 0; i < nn - 1; i++)
        {
            xPoints.push_back(xPoints[i]+ecart);
        }

        //Sort
        std::sort(xPoints.begin(), xPoints.end());
        boost::math::weibull_distribution<> myWeibull(this->shape, this->scale);
        //Calcul les donnees de y
        for(double i = boost::math::quantile(myWeibull,0.0001) ; i < boost::math::quantile(myWeibull,0.9999) ; i += (boost::math::quantile(myWeibull,0.9999) - boost::math::quantile(myWeibull,0.0001))/(double)nn)
        {
            yPoints.push_back(boost::math::pdf(myWeibull,i));
        }

        yPoints.push_front(0);
        yPoints.push_back(0);

        // Maintenant les valeur de l'histogram.

        //grandeur des division
        double ecartItv = this->etendue/nbi;

        vector<double> intervalle;
        for(int i = 0; i < nbi+1; i++)
        {
            intervalle.push_back(donneesIndivitual[0]+i*ecartItv);
        }

        for(int i = 0; i < nbi && i < intervalle.size(); i++)
        {
            pointsX.push_back(intervalle[i]+ecartItv/2);
        }

        double maxAutreY = 0;
        if(ecartItv != 0 && intervalle.size() != 0)
        {
            //hauteur des y
            for(int i = 0; i < nbi; i++)
            {
                //mets tous a zero
                pointsY.push_back(0);
            }

            for(int i = 0; i < n && i < donneesIndivitual.size(); i++)
            {
                int x = ((donneesIndivitual[i] - donneesIndivitual[0]) / ecartItv) - 0.5;
                pointsY[x]++;
                if(pointsY[x] > maxAutreY)
                {
                    maxAutreY = pointsY[x];
                }
            }
        }

    }

    QJsonArray points;
    for(int i=0; i < xPoints.size() && i < yPoints.size(); i++){
        QJsonObject pt;
        pt["q"] = xPoints.at(i);
        pt["p"] = yPoints.at(i);
        points.push_back(pt);
    }
    json["distribution_points"] = points;

    // now the histogram points...
    QJsonArray histogramPoint;
    for(int i=0; i < pointsX.size() && pointsY.size(); i++){
        QJsonObject pt;
        pt["x"] = pointsX.at(i);
        pt["y"] = pointsY.at(i);
        histogramPoint.push_back(pt);
    }

    json["histogram_points"] = histogramPoint;

}
