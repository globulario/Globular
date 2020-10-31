//DonneesAnalyse.h
//Auteur: Tommy Urbain
//Dernière modification:

// Every Data in the sample will be major and convert
// to an integer value.
#ifndef DONNEES_ANALYSE_H
#define DONNEES_ANALYSE_H

#include "SousGroupe.h"
#include "Erreur.h"
#include <vector>
#include <list>
#include <QString>
#include <qvector.h>
#include <sstream>
#include <QJsonObject>

using std::vector;
int getPrecision(double val);

// Les informations relative a une pièce.
struct PieceInfo {
	PieceInfo(){
		isSelect = false;
		tolZone = 0.0f;
		creationDate = "";
		path = "";
		serial = "";
	}

	PieceInfo(const PieceInfo& other){
		isSelect = other.isSelect;
		tolZone = other.tolZone;
		creationDate = other.creationDate;
		path = other.path;
		serial = other.serial;
	}

	// La liste des information du tableau a sauvegardées.
	string serial;
	bool isSelect;
	double tolZone;
	string creationDate;
	double coteZ;
	string path;

};

enum Tendance {
	Convergent,
	Divergent,
	Neutre
};

class DonneesAnalyse
{
public:
	//Constructeur et destructeur
	DonneesAnalyse();
	DonneesAnalyse(const DonneesAnalyse &);
	~DonneesAnalyse();

private:
	vector<PieceInfo*> piecesInfos;

	// Le informations des workpath...
	vector<std::string> unselect_workpaths;

	// Les vecteur qui contienne les information nécessaires a évalué les indices dans le temps...
	vector<double> evolution_moyennes;
	vector<double> evolution_cp;
	vector<double> evolution_cpk;

	//contient la liste des test qui ont échoués...
	QVector<int> failedTests;

	//attributs privé
	double test_K1;
	int test_K2;
	int test_K3;
	int test_K4;
	int test_K5;
	int test_K6;
	int test_K7;
	int test_K8;

	bool state_test_K1;
	bool state_test_K2;
	bool state_test_K3;
	bool state_test_K4;
	bool state_test_K5;
	bool state_test_K6;
	bool state_test_K7;
	bool state_test_K8;

	// Données des limites de contrôles des cartes I/MR
	double LCS_X;
	double LCI_X;
	double LC_X;
	double LCS_MR;
	double LCI_MR;
	double LC_MR;
	
	// Comme les valeurs évolue dans le temps il faut a un moment les figé pour voir
	// comment évolue la capabilité dans le temps.
	double LCX_fix;
	double LCS_X_fix;
	double LCI_X_fix;
	string date_LCX_fix;

	double moyene;
	double ecart_type;
	double etendue;
	double cp;
	double cpk;

	double lotol;
	double uptol;
	double tolNominal;

	double limiteInf;
	double limiteSup;

	string NoTol;
	string NoFeat;
	string OpNo;
	string NoModele;
	string DateDebut;
	string DateFin;
	string Commentaire;
	string tolOption;

	int isNormal;
	double Normality;
	double Skewness;
	double Kurtosis;
	double CV;
	double Variance;
	double Median;
	double AD;
	double Pvalue;
	double Z;
	double shape;
	double scale;

	//Si true, les calculs doivent etre faite pour une population
	//Si false, les calculs doivent etre faite pour un echantillion
	bool PopOuEchan;
	string note;
	string path;


	// Evolution des indice qualité piece par piece...
    void calculEvolution(unsigned int size);

	// Ce champs détermine si l'analyse doit ce mettre automatiquement a jour
	bool isStatic_;

	// Vector temporaire qui garde en mémoire les valeurs.
    std::vector<double> donneesIndividuelle;

public:
    void initToleranceInfo(double tolzon, double lotol, double uptol);

	// Calcul des parametre de la courbe weibull...
    void calculWeibullParameter(double& shape, double& scale, std::vector<double> data,  const double& tolNom);

	// Le chemin du path...
	string getPath();
	void setPath(string);

	// Ajoute une information au sujet d'une piece.
	void clearInfos();
	void addPieceInfo(PieceInfo* info);
    vector<PieceInfo*>& getPieceInfosLst();

	//accesseurs pour les paramètres de tests
	double getTest_K1();
	int getTest_K2();
	int getTest_K3();
	int getTest_K4();
	int getTest_K5();
	int getTest_K6();
	int getTest_K7();
	int getTest_K8();

	//Accesseurs pour les états des tests
	bool getState_test_K1();
	bool getState_test_K2();
	bool getState_test_K3();
	bool getState_test_K4();
	bool getState_test_K5();
	bool getState_test_K6();
	bool getState_test_K7();
	bool getState_test_K8();

	//la liste des tests échoués...
	QVector<int> getFailedTests();
	void setFailedTests(QVector<int>);

	//Mutateurs pour les paramètres de test
	void setTest_K1(double test_K1);
	void setTest_K2(int test_K2);
	void setTest_K3(int test_k3);
	void setTest_K4(int test_K4);
	void setTest_K5(int test_K5);
	void setTest_K6(int test_K6);
	void setTest_K7(int test_K7);
	void setTest_K8(int test_K8);

	//Mutateurs pour les paramètres de test
	void setState_test_K1(bool state_test_K1);
	void setState_test_K2(bool state_test_K2);
	void setState_test_K3(bool state_test_K3);
	void setState_test_K4(bool state_test_K4);
	void setState_test_K5(bool state_test_K5);
	void setState_test_K6(bool state_test_K6);
	void setState_test_K7(bool state_test_K7);
	void setState_test_K8(bool state_test_K8);

	//Accesseur pour les cartes de contrôle I/MR
	double getLCS_X();
	double getLCI_X();
	double getLC_X();

	double getLCS_X_Fix();
	double getLCI_X_Fix();
	double getLC_X_Fix();
	string getLC_X_Date();

	double getLCS_MR();
	double getLCI_MR();
	double getLC_MR();

	double getMoyene();
	double getEtendue();
	double getEcartType();
	double getCp();
	double getCpk();

	double getLoTol();
	double getUpTol();
	double getTolNominal();

	double getLimiteInf();
	double getLimiteSup();

	double getScale();
	double getShape();

	QString getNoTol();
	QString getNoFeat();
	QString getOpNo();
	QString getNoModele();
	QString getDateDebut();
	QString getDateFin();
	QString getCommentaire();
	QString getTolOption();
	QString getDescription();

	int getIsNormal();
	double getNormality();
	double getSkewness();
	double getKurtosis();
	double getCV();
	double getVariance();
	double getMedian();
	double getAD();
	double getPvalue();
	double getZ();
	bool getPopOuEchan();
	string getNote();
    const vector<double>& getDonneesIndividuelle();
    double getDonneesIndividuellePrecision();
    int getDonneesIndividuelleSummation();

	//Matateurs pour les cartes de contrôle I/MR
	void setLCS_X(double LCS_X);
	void setLCI_X(double LCI_X);
	void setLC_X(double LC_X);
	
	void setLCS_X_Fix(double LCS_X);
	void setLCI_X_Fix(double LCI_X);
	void setLC_X_Fix(double LC_X);
	void setLC_X_Date(string date);

	void setLCS_MR(double LCS_MR);
	void setLCI_MR(double LCI_MR);
	void setLC_MR(double LC_MR);

	void setMoyene(double moyene);
	void setEtendue(double etendue);
	void setEcartType(double ecart_type);
	void setCp(double cp);
	void setCpk(double cpk);

	void setLoTol(double lotol);
	void setUpTol(double uptol);
	void setTolNominal(double tolNominal);

	void setLimiteInf(double limiteInf);
	void setLimiteSup(double limiteSup);

	void setNoTol(QString NoTol);
	void setNoFeat(QString NoFeat);
	void setOpNo(QString OpNo);
	void setNoModele(QString NoModele);
	void setDateDebut(QString DateDebut);
	void setDateFin(QString DateFin);
	void setCommentaire(QString Commentaire);
	void setTolOption(QString opTol);

	void setIsNormal(int isNormal);
	void setNormality(double Normality);
	void setSkewness(double Skewness);
	void setKurtosis(double Kurtosis);
	void setCV(double CV);
	void setVariance(double Variance);
	void setMedian(double Median);
	void setAD(double AD);
	void setPvalue(double Pvalue);
	void setZ(double Z);
	void setPopOuEchan(bool PopouEchan);
	void setNote(string& note);
	bool isStatic();
	void setIsStatic(bool);
	void setScale(double _scale);
	void setShape(double _shape);

	// Les tendence...
	Tendance getTendanceMoyenne();
	Tendance getTendanceCp();
	Tendance getTendanceCpk();
	bool isTendanceCpkAcceptable();
	bool isTendanceCpAcceptable();

	bool isSupposedToBeNormal();

	//tableau des sous groupes pour la carte_IMR
	vector<SousGroupe> sousGroupe_IMR;

    //tableau des états de chaque point des 2 courbes pour afficher les point d'erreur en rouge
	vector<Erreur> state_X;
	vector<Erreur> state_MR;

	//tableau contenant les données à être analysé par l'analyseur (graphe Individual value)
    vector<double> donneesIndivitual;

	//tableau contenant les cote z de chaque donnee
	vector<double> coteZDonneeIndivitual;

	//tableau contenant les données à être analysé par l'analyseur (graphe tolerance)
    vector<double> donneesTolerance;

	//surcharge d'opérateur
	DonneesAnalyse& operator =(const DonneesAnalyse &donnees);

	//vide le tableau des pieces..
	void clearPieces();

	//méthode pour trier les pièces par date..
	void sortPieceByDate();

	// Mise a jour des donnée calculé de l'analyse.
	void update();

	// Mise a jour d'un info de la liste des information...
    void update(std::string serial, bool isSelect);

	// Cette fonction permet d'initialiser la liste des info dans la classe... 
	// * Les donnée relative au modèle, numéros d'opération etc, doivent avoir été initialisé correctement
	//   avant l'appel de cette méthode...
	void initItemInfo();

	// Permet d'ajouter une route a soustraire de l'analyse...
	void addUnselectWorkPath(QString);
	void removeUnselectWorkPaht(QString);
	void setUnselectWorkPath(QList<QString>& lst);
	QList<QString>* getUnselectWorkPath();

    // Convertion to and from json values.
    void read(const QJsonObject &json);
    void write(QJsonObject &json) const;

};

#endif
