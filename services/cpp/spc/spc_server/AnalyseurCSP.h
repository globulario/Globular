//AnalyseurCSP.h
//Auteur: Tommy Urbain
//Modifier par: Antoine Mathieu Morel
//Dernière modification: 11 juillet 2011
//Classe permettant de faire l'analyse des données

#ifndef ANALYSEURCSP_H
#define ANALYSEURCSP_H

#include "SousGroupe.h"
#include "DonneesAnalyse.h"
#include "Erreur.h"
#include <vector>
#include <QString>
#include <algorithm>
#include <cmath>
#include <iostream>
#include <fstream>

using namespace std;

//*************************************
#include <boost/math/distributions/normal.hpp> // for normal_distribution
  using boost::math::normal; // typedef provides default type is double.
#include <boost/math/distributions/cauchy.hpp> // for cauchy_distribution
  using boost::math::cauchy; // typedef provides default type is double.
#include <boost/math/distributions/find_location.hpp>
  using boost::math::find_location;
#include <boost/math/distributions/find_scale.hpp>
  using boost::math::find_scale;
  using boost::math::complement;
  using boost::math::policies::policy;
//********************************

#define PI 3.14159265358979323846
#define ENUM 2.71828182845904523536

class AnalyseurCSP
{
private:
	//Attributs privées
	double range_I; // C'est la différence entre la plus grande valeur et la plus petite
	double moyenne_I;// Moyenne des données à être évaluées
	double moyenne_moyenne_I; //Moyenne de la mmoyenne des sous-groupe
	double range_R; // C'est la différence entre la plus grande valeur et la plus petite
	double moyenne_R;// Moyenne des données sous groupes
	double moyenneRange_R; //La moyenne des écarts
	double moyenneRange_IMR;

	//Tableau de constante pour l'analyse
	std::vector<double> A2;
	std::vector<double> D3;
	std::vector<double> D4;
	std::vector<double> A3;
	std::vector<double> B3;
	std::vector<double> B4;

	bool threadGoing;

public:
	//Constructeur et Destructeur
	AnalyseurCSP();
	~AnalyseurCSP();

	bool getThreadGoing();


	//Autres Fonctions
	void calculIMR(DonneesAnalyse &donnees);
	void verifierRegleStatistique(DonneesAnalyse &donnees);
	void getDonnees(DonneesAnalyse &donnees);
	void analyseComplete(DonneesAnalyse &donnees);
	void ajouterFailedTests(DonneesAnalyse &donnees, int noErreur);

	void calculMoyenne(DonneesAnalyse &donnees);
	void calculEtendue(DonneesAnalyse &donnees);
	void calculEcartType(DonneesAnalyse &donnees);
	void calculVariance(DonneesAnalyse &donnees);
	void calculMedian(DonneesAnalyse &donnees);
	void calculCV(DonneesAnalyse &donnees);
	void calculCpCpk(DonneesAnalyse &donnees);
	void calculSkewness(DonneesAnalyse &donnees);
	void calculKurtosis(DonneesAnalyse &donnees);
	void calculCoteZ(DonneesAnalyse &donnees);
	void calculZ(DonneesAnalyse &donnees);

	void getNormalInfo(DonneesAnalyse &donnees);
	bool isNormalPro(DonneesAnalyse &donnees);
	bool AndersonDarlingTest(DonneesAnalyse &donnees);
	int SkewnessTest(DonneesAnalyse &donnees);
	int KurtosisTest(DonneesAnalyse &donnees);
	bool SkewnessKurtosisAll(DonneesAnalyse &donnees);
	double KolSmiTest(DonneesAnalyse &donnees);
	double randn_notrig(double mu, double sigma);
	double randn_trig(double mu, double sigma);
	double KolmogorovTest(int na, vector<double> a, int nb, vector<double> b);
	double KolmogorovProb(double z);
	int Nint(float x);
	int Nint(double x);
	double AndersonDarlingPValue(int n, double z);
	double AndersonDarlingStatistic(DonneesAnalyse &donnees);
	double CDF(double x, double mu, double sigma);
	void BoxCoxTransformation(DonneesAnalyse &donnees, vector<DonneesAnalyse*> &output);
};
#endif
