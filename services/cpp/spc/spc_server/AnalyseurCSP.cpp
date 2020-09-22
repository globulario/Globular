//AnalyseurCSP.cpp
//fait le 11 mars 2011

#include "AnalyseurCSP.h"
#include <boost/math/distributions/weibull.hpp>
#include <boost/multiprecision/cpp_int.hpp>
#include <boost/multiprecision/gmp.hpp>
using namespace boost::multiprecision;

#include <sstream>

#include <qdebug.h>
#if defined(WIN32) && defined(_DEBUG)
#define DEBUG_NEW new( _NORMAL_BLOCK, __FILE__, __LINE__ )
#define new DEBUG_NEW
#endif

#define MAJORING_FACTOR 100000 // Cette valeur permet de transformer un double en nombre entier...

using namespace std;

double round_to_digits(double f, int n){
    int p = pow(10, n);
    f *= p;

    if(f > 0){
        f += 0.5;
    }else if(f < 0){
        f -= 0.5;
    }

    f = (float)((int) f);
    if(p > 0){
        f /= p;
    }
    return f;
}

//Constructeur et Destructeur
AnalyseurCSP::AnalyseurCSP()
{
    //Initialisation des tableaux de constante
    A3.push_back(NULL); //0 n'existe pas
    A3.push_back(NULL); //1 n'existe pas
    A3.push_back(2.66); //2
    A3.push_back(1.95); //3
    A3.push_back(1.63); //4
    A3.push_back(1.43); //5
    A3.push_back(1.29); //6
    A3.push_back(1.18); //7
    A3.push_back(1.1);  //8
    A3.push_back(1.03); //9
    A3.push_back(0.98); //10

    B3.push_back(NULL); //0 n'existe pas
    B3.push_back(NULL); //1 n'existe pas
    B3.push_back(0);    //2
    B3.push_back(0);    //3
    B3.push_back(0);    //4
    B3.push_back(0);    //5
    B3.push_back(0.03); //6
    B3.push_back(0.12); //7
    B3.push_back(0.19); //8
    B3.push_back(0.24); //9
    B3.push_back(0.28); //10

    B4.push_back(NULL); //0 n'existe pas
    B4.push_back(NULL); //1 n'existe pas
    B4.push_back(3.27); //2
    B4.push_back(2.57); //3
    B4.push_back(2.27); //4
    B4.push_back(2.09); //5
    B4.push_back(1.97); //6
    B4.push_back(1.88); //7
    B4.push_back(1.82); //8
    B4.push_back(1.76); //9
    B4.push_back(1.72); //10

    A2.push_back(NULL); //0 n'existe pas
    A2.push_back(NULL); //1 n'existe pas
    A2.push_back(1.88); //2
    A2.push_back(1.02); //3
    A2.push_back(0.73); //4
    A2.push_back(0.58); //5
    A2.push_back(0.48); //6
    A2.push_back(0.42); //7
    A2.push_back(0.37); //8
    A2.push_back(0.34); //9
    A2.push_back(0.31); //10

    D3.push_back(NULL); //0
    D3.push_back(NULL); //1
    D3.push_back(0); //2
    D3.push_back(0); //3
    D3.push_back(0); //4
    D3.push_back(0); //5
    D3.push_back(0); //6
    D3.push_back(0.08); //7
    D3.push_back(0.14); //8
    D3.push_back(0.18); //9
    D3.push_back(0.22); //10

    D4.push_back(NULL); //0
    D4.push_back(NULL); //1
    D4.push_back(3.27); //2
    D4.push_back(2.57); //3
    D4.push_back(2.28); //4
    D4.push_back(2.11); //5
    D4.push_back(2); //6
    D4.push_back(1.92); //7
    D4.push_back(1.86); //8
    D4.push_back(1.82); //9
    D4.push_back(1.78); //10

    threadGoing = false;
}
AnalyseurCSP::~AnalyseurCSP()
{
}

//Autres Fonctions


//Cette class doit avoir NoTol, OpNo, NoModele, DateDebut, DateFin deja avec des donnees dans l'objet donnees DonneesAnalyse
void AnalyseurCSP::analyseComplete(DonneesAnalyse &donnees)
{
    donnees.initItemInfo();
    if(donnees.getDonneesIndividuelle().size() > 0)
    {
        calculIMR(donnees);
    }

}

//Permet de calculer les valeurs pour les limites de contrôle LC, LCS et LCI pour les cartes I/MR
//Intrant: un tableau de double contenant les données à être analysées
//Extrant: aucun (les données membres concernant les valeurs des cartes I/MR sont initialisées
void AnalyseurCSP::calculIMR(DonneesAnalyse &donnees)
{

    //Constante pour calculer les limites de contrôle
    const double d2 = 1.128;
    const double D3 = this->D3[2];
    const double D4 = this->D4[2];
    donnees.sousGroupe_IMR.clear();
    const vector<double>& donneesIndividuelle = donnees.getDonneesIndividuelle();

    SousGroupe s;
    std::vector<double>::const_iterator it = donneesIndividuelle.begin();
    if(donnees.getDonneesIndividuelle().size() > 1){
        while(it != donneesIndividuelle.end() - 1)
        {
            s.donnees.push_back(*it);
            s.donnees.push_back(*(it + 1));
            it++;
            donnees.sousGroupe_IMR.push_back(s);
            s.donnees.clear();
        }

    }else if(donneesIndividuelle.size() == 1){
        s.donnees.push_back(*it);
        s.donnees.push_back(*it);
        donnees.sousGroupe_IMR.push_back(s);
    }

    // Calcul des moyenne ecart type et variance.
    calculMoyenne(donnees);
    calculVariance(donnees);
    calculEcartType(donnees);

    for(int i = 0; i < donnees.sousGroupe_IMR.size(); ++i)
    {
        donnees.sousGroupe_IMR[i].setRange(abs(donnees.sousGroupe_IMR[i].donnees[0] - donnees.sousGroupe_IMR[i].donnees[1]));
    }

    moyenneRange_IMR = 0;

    for(int i = 0; i < donnees.sousGroupe_IMR.size(); ++i)
    {
        this->moyenneRange_IMR += donnees.sousGroupe_IMR[i].getRange();
    }

    this->moyenneRange_IMR = (moyenneRange_IMR / donnees.sousGroupe_IMR.size());

    // utilise la valeur de la moyenne...
    double valeurCentrale = donnees.getMoyene();

    //Initialisation des courbes des limites de contrôle LC, LCS et LCI
    donnees.setLC_MR(moyenneRange_IMR);
    donnees.setLCS_MR(D4 * moyenneRange_IMR);
    donnees.setLCI_MR(D3 * moyenneRange_IMR);

    //Initialisation des courbes de limites de contrôle LC_I, LCS_I et LCI_I
    donnees.setLC_X(valeurCentrale);
    donnees.setLCS_X(valeurCentrale +  ((3/d2) * moyenneRange_IMR));
    donnees.setLCI_X(valeurCentrale - ((3/d2) * moyenneRange_IMR));

    calculEtendue(donnees);
    calculMedian(donnees);
    calculCV(donnees);
    calculSkewness(donnees);
    calculKurtosis(donnees);
    calculCoteZ(donnees);

    donnees.setAD(AndersonDarlingStatistic(donnees));
    donnees.setPvalue(AndersonDarlingPValue(int(donneesIndividuelle.size()),donnees.getAD()));

    getNormalInfo(donnees);

    calculCpCpk(donnees);

    calculZ(donnees);

    verifierRegleStatistique(donnees);
}

//Calcul la Moyenne des donnees
void AnalyseurCSP::calculMoyenne(DonneesAnalyse &donnees)
{
    // Calcul de la moyenne...
    double moyenneX = 0.0;
    for (int i = 0; i < donnees.getDonneesIndividuelle().size(); ++i) {
        moyenneX += donnees.getDonneesIndividuelle()[i];
    }
    donnees.setMoyene(moyenneX / donnees.getDonneesIndividuelle().size());
}

//Calcul la Etendue des donnees
void AnalyseurCSP::calculEtendue(DonneesAnalyse &donnees)
{
    if(donnees.getDonneesIndividuelle().size() > 0)
    {
        vector<double> donneesIndivitual2(donnees.getDonneesIndividuelle().size());

        copy(donnees.getDonneesIndividuelle().begin(), donnees.getDonneesIndividuelle().end(), donneesIndivitual2.begin());

        std::sort(donneesIndivitual2.begin(), donneesIndivitual2.end());

        donnees.setEtendue(donneesIndivitual2[donneesIndivitual2.size()-1] - donneesIndivitual2[0]);
    }
}

//Calcul l'EcartType des donnees
void AnalyseurCSP::calculEcartType(DonneesAnalyse &donnees)
{
    double ecartType = 0;
    ecartType = sqrt(donnees.getVariance());
    donnees.setEcartType(ecartType);
}

int getValue(double val, int precision){
    std::stringstream ss;
    ss << val;
    string strVal = ss.str();
    size_t start = strVal.find(".");

    std::string major;
    std::string minor;
    if(start != std::string::npos){
        major = strVal.substr(0, start);
        minor = strVal.substr(start + 1);
    }else{
        major = strVal;
        minor = "0";
    }


    // Dans le cas ou la precision n'est pas atteinte
    while(minor.length() < precision){
        minor += "0";
    }

    // Dans le cas ou elle est dépassée...
    if(minor.length() > precision){
        minor = minor.substr(0, precision);
    }

    strVal = major + minor;
    int intVal = atoi(strVal.c_str());

    return intVal;
}

//Calcul la Variance des donnees
void AnalyseurCSP::calculVariance(DonneesAnalyse &donnees)
{
    double variance = 0;
    if(donnees.getDonneesIndividuelle().size()  <= 1){
        donnees.setVariance(variance);
        return;
    }

    int precision = donnees.getDonneesIndividuellePrecision();
    double factor = pow(10, precision);
    double nomVal = donnees.getTolNominal();
    int nominal = getValue(nomVal, precision);

    mpf_float_100 sum(0);//int sum = 0;

    size_t size = donnees.getDonneesIndividuelle().size();
    int shift = 10000;
    for (int i = 0; i < size; ++i) {
        sum += shift * getValue(donnees.getDonneesIndividuelle()[i], precision) - nominal; // Le resultat peu contenir une valeur signification
    }

    mpf_float_100 moyenne = sum / size;

    // Je vais réutilisé la variable sum dans le calcul de la variance...
    sum = 0;
    for (int i = 0; i < size; ++i) {
        int tmp = shift * getValue(donnees.getDonneesIndividuelle()[i], precision) - nominal;
        sum += pow( tmp - moyenne, 2);
    }

    //qDebug() << variance;
    if(donnees.getPopOuEchan())
    {
        mpf_float_100 denominator(size * pow(factor * shift, 2));
        mpf_float_100 v = sum / denominator;
        mpf_t f;
        mpf_init(f);
        mpf_set(f, v.backend().data());
        variance = v.convert_to<double>();
    }
    else if(!donnees.getPopOuEchan())
    {
        mpf_float_100 denominator((size-1) * pow(factor * shift, 2));
        mpf_float_100 v = sum / denominator;
        mpf_t f;
        mpf_init(f);
        mpf_set(f, v.backend().data());
        variance = v.convert_to<double>();
    }
    else
    {
        variance = 0;
    }

    donnees.setVariance(variance);
}

//Calcul la Median des donnees
void AnalyseurCSP::calculMedian(DonneesAnalyse &donnees)
{

    double median = 0;
    if(donnees.getDonneesIndividuelle().size() >= 1)
    {
        if (donnees.getDonneesIndividuelle().size()%2 == 1) {
            //        if odd number of elements, median is in the middle of the sorted array
            median = donnees.getDonneesIndividuelle()[donnees.getDonneesIndividuelle().size()/2];
        } else {
            //        if even number of elements, median is avg of two elements straddling
            //        the middle of the array
            median = (donnees.getDonneesIndividuelle()[donnees.getDonneesIndividuelle().size()/2] + donnees.getDonneesIndividuelle()[(donnees.getDonneesIndividuelle().size()/2)-1])/2;
        }
    }
    donnees.setMedian(median);
}

//Calcul le Coefficient de variation
void AnalyseurCSP::calculCV(DonneesAnalyse &donnees)
{
    double cv;
    if (donnees.getMoyene() == 0)
    {
        cv = 0;
    }
    else
    {
        cv = donnees.getEcartType() / abs(donnees.getMoyene());
    }
    donnees.setCV(cv);
}

//Calcul le Cp et Cpk des donnees
void AnalyseurCSP::calculCpCpk(DonneesAnalyse &donnees)
{
    donnees.setCp(0);
    donnees.setCpk(0);

    if(donnees.isSupposedToBeNormal()){

        if (donnees.getEcartType() == 0)
        {
            return;
        }

        const double NB_ECART_TYPE = 6;
        const double MOITIER_NB_ECART_TYPE = 3;
        double cp = (donnees.getLimiteSup() - donnees.getLimiteInf()) / (NB_ECART_TYPE * donnees.getEcartType());
        donnees.setCp(cp);

        double Zmin = ((donnees.getMoyene() - donnees.getLimiteInf()) / (MOITIER_NB_ECART_TYPE * donnees.getEcartType()));
        double Zmax = ((donnees.getLimiteSup() - donnees.getMoyene()) / (MOITIER_NB_ECART_TYPE * donnees.getEcartType()));

        if(Zmin < Zmax)
            donnees.setCpk(Zmin);
        else
            donnees.setCpk(Zmax);
    }else{
        if (donnees.getScale() == 0 || donnees.getShape() == 0)
        {
            return;
        }

        // Calcul le ppk...
        boost::math::weibull_distribution<> myWeibull(donnees.getShape(), donnees.getScale());
        double alpha = 0.00135;
        double middleVal = boost::math::quantile(myWeibull,0.5);	// 0.5
        double rightVal = boost::math::quantile(myWeibull,1-alpha);	// 1-alpha
        double Zmax = (donnees.getLimiteSup() - middleVal)/(rightVal - middleVal);

        donnees.setCpk(Zmax);
        donnees.setCp(-1);
    }
}

//Calcul le Skewness des donnees, aussi nommer le Coefficient d'asymetrie
void AnalyseurCSP::calculSkewness(DonneesAnalyse &donnees)
{
    //http://www.indiana.edu/~educy520/sec5982/week_12/skewness_demo.pdf
    //http://www.dreamincode.net/code/snippet1447.htm
    //Formule utuliser viens de l'aide de excel
    //http://www.mathworks.com/help/toolbox/stats/skewness.html
    double skewness = 0;
    double n = donnees.getDonneesIndividuelle().size();
    double somme1 = 0;
    double somme2 = 0;

    if (n <= 1){
        donnees.setSkewness(0);
        return;
    }

    for(int i = 0; i < n; ++i)
    {
        somme1 += pow( donnees.getDonneesIndividuelle()[i]-donnees.getMoyene() , 3.0);
        somme2 += pow( donnees.getDonneesIndividuelle()[i]-donnees.getMoyene() , 2.0);
    }

    if (pow( sqrt(somme2/n) ,3.0) == 0)
    {
        skewness = 0;
    }
    else
    {
        skewness = (somme1/n)/(pow( sqrt(somme2/n) ,3.0));
    }

    if(donnees.getPopOuEchan())
    {
        skewness = (sqrt(n*(n-1))* skewness)/(n-2);
    }

    donnees.setSkewness(skewness);
}

//Calcul le Kurtosis des donnees, aussi nommer le Coefficient d'aplatissement
void AnalyseurCSP::calculKurtosis(DonneesAnalyse &donnees)
{
    // http://www.dreamincode.net/code/snippet1447.htm
    // Formule utuliser viens de l'aide de excel
    // http://www.mathworks.com/help/toolbox/stats/kurtosis.html
    double kustosis = 0;
    double n = donnees.getDonneesIndividuelle().size();
    double somme1 = 0;
    double somme2 = 0;

    if(n <= 1){
        donnees.setKurtosis(0);
        return;
    }

    for(int i = 0; i < n; ++i)
    {
        somme1 += pow( (donnees.getDonneesIndividuelle()[i]-donnees.getMoyene()) , 4.0);
        somme2 += pow( (donnees.getDonneesIndividuelle()[i]-donnees.getMoyene()) , 2.0);
    }

    if (pow( (somme2/n) , 2.0) == 0)
    {
        kustosis = 0;
    }
    else
    {
        kustosis = (somme1/n) / pow( (somme2/n) , 2.0);
    }

    if(donnees.getPopOuEchan())
    {
        kustosis = ( (n-1)*( (n+1)*kustosis-3*(n-1)  ) ) / ( (n-2)*(n-3) ) + 3;
    }

    //Par ce que habituellement une normal a un kurtosis de 3, mais en faisant moins 3
    //ont dites que le kurtosis d'une normal est zero. C'est seulment une questions de comment
    //l'utilisateur veut interpreter. Les fonction "KurtosisTest" et "SkewnessKurtosisAll"
    //prends cette conditions en consideration pour obtenir le resultat
    kustosis = kustosis - 3;
    donnees.setKurtosis(kustosis);

}

//Calcul les cote Z pour chaque donnees
void AnalyseurCSP::calculCoteZ(DonneesAnalyse &donnees)
{
    donnees.coteZDonneeIndivitual.clear();
    for(int i = 0; i < donnees.getDonneesIndividuelle().size(); i++)
    {
        if(donnees.getEcartType() == 0){
            donnees.coteZDonneeIndivitual.push_back(0);
        }else{
            donnees.coteZDonneeIndivitual.push_back((donnees.getDonneesIndividuelle().at(i) - donnees.getMoyene())/donnees.getEcartType());
        }
    }
}

//Calcul de Z
void AnalyseurCSP::calculZ(DonneesAnalyse &donnees)
{
    double zMax = 0;
    double zMin = 0;

    if(donnees.getDonneesIndividuelle().size() > 1)
    {
        if(donnees.getUpTol() != 0)
        {
            zMax = 1 - CDF(donnees.getTolNominal() + donnees.getUpTol(), donnees.getMoyene(), donnees.getEcartType());
            zMax *= 100;
        }

        if(donnees.getLoTol() != 0)
        {
            zMin = CDF(donnees.getTolNominal() + donnees.getLoTol(), donnees.getMoyene(), donnees.getEcartType());
            zMin *= 100;
        }

        donnees.setZ(zMax+zMin);
    }
}

//Permet de vérfifer s'il y a des données ayant des comportements suspets
//Intrant: Un tableau contemant les données
//Extrant: Un tableau d'état est créer avec les status de cchaque points
//NOTE://pour pouvoir utiliser cette fonction, il faut trouver le moyen de générer deux types de symboles pour la même courbe.
void AnalyseurCSP::verifierRegleStatistique(DonneesAnalyse &donnees)
{
    //initialisation de l'états de chaques points individual
    for(int i = 0 ; i < donnees.getDonneesIndividuelle().size(); ++i)
    {
        Erreur err;
        donnees.state_X.push_back(err);
    }

    //initialisation de l'états de chaques points moving rage
    for(int i = 0 ; i < donnees.sousGroupe_IMR.size(); ++i)
    {
        Erreur err;
        donnees.state_MR.push_back(err);
    }

    //1- 1 point dans la zone hors contrôle(1 point > ou < que +-3 écart types de la moyenne
    if(donnees.getState_test_K1())
    {
        //Pour le graphique individual
        for(int i = 0 ; i < donnees.getDonneesIndividuelle().size(); ++i)
        {
            if(donnees.getDonneesIndividuelle()[i] > donnees.getLCS_X() || donnees.getDonneesIndividuelle()[i] < donnees.getLCI_X() )
            {
                donnees.state_X[i].noErreurs.push_back(1);
                ajouterFailedTests(donnees,1);
            }
        }

        //Pour le graphique moving range
        for(int i = 0 ; i < donnees.sousGroupe_IMR.size(); ++i)
        {
            if(donnees.sousGroupe_IMR[i].getRange() > donnees.getLCS_MR() || donnees.sousGroupe_IMR[i].getRange() < donnees.getLCI_MR())
            {
                donnees.state_MR[i].noErreurs.push_back(1);
                //ajouterFailedTests(donnees,1);
            }
        }
    }

    //2- 6 points consécutifs dans la même zone
    if(donnees.getState_test_K2())
    {
        //Pour le graphique individual
        if(donnees.getDonneesIndividuelle().size() > donnees.getTest_K2())
        {
            int nbSuiventHaut = 0;
            int nbSuiventBas = 0;

            for(int i = 1; i <  donnees.getDonneesIndividuelle().size(); ++i)
            {
                if(donnees.getDonneesIndividuelle()[i] > donnees.getLC_X())
                {
                    nbSuiventHaut++;
                }
                else
                {
                    nbSuiventHaut = 0;
                }

                if(donnees.getDonneesIndividuelle()[i] < donnees.getLC_X())
                {
                    nbSuiventBas++;
                }
                else
                {
                    nbSuiventBas = 0;
                }

                if(nbSuiventHaut >= donnees.getTest_K2() || nbSuiventBas >= donnees.getTest_K2())
                {
                    donnees.state_X[i].noErreurs.push_back(2);
                    ajouterFailedTests(donnees,2);
                }
            }
        }

        //Pour le graphique moving range
        if(donnees.sousGroupe_IMR.size() > donnees.getTest_K2())
        {
            int nbSuiventHaut = 0;
            int nbSuiventBas = 0;

            for(int i = 1; i <  donnees.sousGroupe_IMR.size(); ++i)
            {
                if(donnees.sousGroupe_IMR[i].getRange() > donnees.getLC_MR())
                {
                    nbSuiventHaut++;
                }
                else
                {
                    nbSuiventHaut = 0;
                }
                if(donnees.sousGroupe_IMR[i].getRange() < donnees.getLC_MR())
                {
                    nbSuiventBas++;
                }
                else
                {
                    nbSuiventBas = 0;
                }
                if(nbSuiventHaut >= donnees.getTest_K2() || nbSuiventBas >= donnees.getTest_K2())
                {
                    donnees.state_MR[i].noErreurs.push_back(2);
                    //ajouterFailedTests(donnees,2);
                }
            }
        }
    }

    //3- 6 point augmentant ou diminuant consécutivement(dérives)
    if(donnees.getState_test_K3())
    {
        //pour le graphique individual
        if(donnees.getDonneesIndividuelle().size() > donnees.getTest_K3())
        {
            int donneesAugmente = 0;
            int donneesDiminue = 0;

            for(int i = 1; i < donnees.getDonneesIndividuelle().size(); ++i)
            {
                if(donnees.getDonneesIndividuelle()[i] > donnees.getDonneesIndividuelle()[i-1])
                {
                    donneesAugmente++;
                }
                else
                {
                    donneesAugmente = 0;
                }
                if(donnees.getDonneesIndividuelle()[i] < donnees.getDonneesIndividuelle()[i-1])
                {
                    donneesDiminue++;
                }
                else
                {
                    donneesDiminue = 0;
                }
                if(donneesDiminue >= donnees.getTest_K3() || donneesAugmente >= donnees.getTest_K3())
                {
                    donnees.state_X[i].noErreurs.push_back(3);
                    ajouterFailedTests(donnees,3);
                }
            }
        }

        //pour le graphique moving range
        if(donnees.sousGroupe_IMR.size() > donnees.getTest_K3())
        {
            int donneesAugmente = 0;
            int donneesDiminue = 0;

            for(int i = 1; i < donnees.sousGroupe_IMR.size(); ++i)
            {
                if(donnees.sousGroupe_IMR[i].getRange() > donnees.sousGroupe_IMR[i-1].getRange())
                {
                    donneesAugmente++;
                }
                else
                {
                    donneesAugmente = 0;
                }
                if(donnees.sousGroupe_IMR[i].getRange() < donnees.sousGroupe_IMR[i-1].getRange())
                {
                    donneesDiminue++;
                }
                else
                {
                    donneesDiminue = 0;
                }
                if(donneesDiminue >= donnees.getTest_K3() || donneesAugmente >= donnees.getTest_K3())
                {
                    donnees.state_MR[i].noErreurs.push_back(3);
                    //ajouterFailedTests(donnees,3);
                }
            }
        }
    }

    //4- Altération autour d'une valeur autre que la moyenne (biais)  --14
    if(donnees.getState_test_K4())
    {
        //pour le graphique individual
        if(donnees.getDonneesIndividuelle().size() > donnees.getTest_K4())
        {
            bool croissant = NULL;
            int nbOssilation = 0;

            for(int i = 1; i < donnees.getDonneesIndividuelle().size(); ++i)
            {
                bool actuelleCroissant;

                if(donnees.getDonneesIndividuelle()[i] > donnees.getDonneesIndividuelle()[i-1])
                {
                    actuelleCroissant = true;
                }
                else
                {
                    actuelleCroissant = false;
                }
                if(croissant != actuelleCroissant)
                {
                    nbOssilation++;
                }
                else
                {
                    nbOssilation = 0;
                }

                croissant = actuelleCroissant;

                if(nbOssilation >= donnees.getTest_K4())
                {
                    donnees.state_X[i].noErreurs.push_back(4);
                    ajouterFailedTests(donnees,4);
                }
            }
        }

        //Pour le graphique moving range
        if(donnees.sousGroupe_IMR.size() > donnees.getTest_K4())
        {
            bool croissant = NULL;
            int nbOssilation = 0;

            for(int i = 1; i < donnees.sousGroupe_IMR.size(); ++i)
            {
                bool actuelleCroissant;

                if(donnees.sousGroupe_IMR[i].getRange() > donnees.sousGroupe_IMR[i-1].getRange())
                {
                    actuelleCroissant = true;
                }
                else
                {
                    actuelleCroissant = false;
                }
                if(croissant != actuelleCroissant)
                {
                    nbOssilation++;
                }
                else
                {
                    nbOssilation = 0;
                }

                croissant = actuelleCroissant;

                if(nbOssilation >= donnees.getTest_K4())
                {
                    donnees.state_MR[i].noErreurs.push_back(4);
                    //ajouterFailedTests(donnees,4);
                }
            }
        }
    }
    //5- K points sur k+1 sont dans la zone A  --2
    if(donnees.getState_test_K5())
    {
        //Pour le graphique individual
        if(donnees.getDonneesIndividuelle().size() > donnees.getTest_K5() + 1)
        {
            for(int i = 0; i < donnees.getDonneesIndividuelle().size(); ++i)
            {
                int nbdansZoneUp = 0;
                int nbdansZoneDown = 0;

                for(int j = i; j < (i+donnees.getTest_K5()+1) && j < donnees.getDonneesIndividuelle().size(); ++j)
                {
                    if(donnees.getDonneesIndividuelle()[j] <= donnees.getLCS_X() &&
                            donnees.getDonneesIndividuelle()[j] >= ((donnees.getLCS_X()-donnees.getLC_X())*(2.0/3.0))+donnees.getLC_X())
                    {
                        nbdansZoneUp++;
                    }

                    if(nbdansZoneUp >= donnees.getTest_K5())
                    {
                        donnees.state_X[j].noErreurs.push_back(5);
                        ajouterFailedTests(donnees,5);
                    }
                }

                for(int j = i; j < (i+donnees.getTest_K5()+1) && j < donnees.getDonneesIndividuelle().size(); ++j)
                {
                    if(donnees.getDonneesIndividuelle()[j] >= donnees.getLCI_X() &&
                            donnees.getDonneesIndividuelle()[j] <= ((donnees.getLCI_X()-donnees.getLC_X())*(2.0/3.0))+donnees.getLC_X())
                    {
                        nbdansZoneDown++;
                    }

                    if(nbdansZoneDown >= donnees.getTest_K5())
                    {
                        donnees.state_X[j].noErreurs.push_back(5);
                        ajouterFailedTests(donnees,5);
                    }
                }
            }

        }

        //Pour le graphique moving range
        if(donnees.sousGroupe_IMR.size() > donnees.getTest_K5()+1)
        {
            for(int i = 0; i <donnees.sousGroupe_IMR.size(); ++i)
            {
                int nbdansZoneUp = 0;
                int nbdansZoneDown = 0;

                for(int j = i; j < (i+donnees.getTest_K5()+1) && j < donnees.sousGroupe_IMR.size();++j)
                {
                    if(donnees.sousGroupe_IMR[j].getRange() <= donnees.getLCS_MR() &&
                            donnees.sousGroupe_IMR[j].getRange() >= ((donnees.getLCS_MR()-donnees.getLC_MR())*(2.0/3.0))+donnees.getLC_MR())
                    {
                        nbdansZoneUp++;
                    }

                    if(nbdansZoneUp >= donnees.getTest_K5())
                    {
                        donnees.state_MR[j].noErreurs.push_back(5);
                        //ajouterFailedTests(donnees,5);
                    }
                }

                for(int j = i; j < (i+donnees.getTest_K5()+1) && j < donnees.sousGroupe_IMR.size();++j)
                {
                    if(donnees.sousGroupe_IMR[j].getRange() >= donnees.getLCI_MR() &&
                            donnees.sousGroupe_IMR[j].getRange() <= ((donnees.getLCI_MR()-donnees.getLC_MR())*(2.0/3.0))+donnees.getLC_MR())
                    {
                        nbdansZoneDown++;
                    }

                    if(nbdansZoneDown >= donnees.getTest_K5())
                    {
                        donnees.state_MR[j].noErreurs.push_back(5);
                        //ajouterFailedTests(donnees,5);
                    }
                }
            }
        }
    }

    //6- quatre point sur cing sont à l'extérieur de la zone C  --4
    if(donnees.getState_test_K6())
    {
        //Pour le graphique individual
        if(donnees.getDonneesIndividuelle().size() > donnees.getTest_K6()+1)
        {
            for(int i = 0 ; i < donnees.getDonneesIndividuelle().size();++i)
            {
                int nbdansZoneUp = 0;
                int nbdansZoneDown = 0;

                for(int j = i; j < (i+donnees.getTest_K6()+1) && j < donnees.getDonneesIndividuelle().size();++j)
                {
                    if(donnees.getDonneesIndividuelle()[j] >= ((donnees.getLCS_X()-donnees.getLC_X())*(1.0/3.0))+donnees.getLC_X())
                    {
                        nbdansZoneUp++;
                    }

                    if(nbdansZoneUp == donnees.getTest_K6())
                    {
                        donnees.state_X[j].noErreurs.push_back(6);
                        ajouterFailedTests(donnees,6);
                    }
                }
                for(int j = i; j < (i+donnees.getTest_K6()+1) && j < donnees.getDonneesIndividuelle().size();++j)
                {
                    if(donnees.getDonneesIndividuelle()[j] <= ((donnees.getLCI_X()-donnees.getLC_X())*(1.0/3.0))+donnees.getLC_X())
                    {
                        nbdansZoneDown++;
                    }

                    if(nbdansZoneDown == donnees.getTest_K6())
                    {
                        donnees.state_X[j].noErreurs.push_back(6);
                        ajouterFailedTests(donnees,6);
                    }
                }
            }
        }
        //pour le graphique moving range
        if(donnees.sousGroupe_IMR.size() > donnees.getTest_K6()+1)
        {
            for(int i = 0; i < donnees.sousGroupe_IMR.size(); ++i)
            {
                int nbdansZoneUp = 0;
                int nbdansZoneDown = 0;

                for(int j = i; j < (i+donnees.getTest_K6()+1) && j < donnees.sousGroupe_IMR.size();++j)
                {
                    if(donnees.sousGroupe_IMR[j].getRange() >= ((donnees.getLCS_MR()-donnees.getLC_MR())*(1.0/3.0))+donnees.getLC_MR())
                    {
                        nbdansZoneUp++;
                    }

                    if(nbdansZoneUp == donnees.getTest_K6())
                    {
                        donnees.state_MR[j].noErreurs.push_back(6);
                        //ajouterFailedTests(donnees,6);
                    }
                }
                for(int j = i; j < (i+donnees.getTest_K6()+1) && j < donnees.sousGroupe_IMR.size();++j)
                {
                    if(donnees.sousGroupe_IMR[j].getRange() <= ((donnees.getLCI_MR()-donnees.getLC_MR())*(1.0/3.0))+donnees.getLC_MR())
                    {
                        nbdansZoneDown++;
                    }

                    if(nbdansZoneDown == donnees.getTest_K6())
                    {
                        donnees.state_MR[j].noErreurs.push_back(6);
                        //ajouterFailedTests(donnees,6);
                    }
                }
            }

        }
    }

    //7- Quinze points consécutifs dans les zone C autour de la ligne de contrôle  --15
    if(donnees.getState_test_K7())
    {
        //Pour le graphique individual
        if(donnees.getDonneesIndividuelle().size() > donnees.getTest_K7())
        {
            int nbPointConsecutif = 0;

            for(int i=0;i < donnees.getDonneesIndividuelle().size();++i)
            {
                if(donnees.getDonneesIndividuelle()[i] > donnees.getLC_X() &&
                        donnees.getDonneesIndividuelle()[i] < ((donnees.getLCS_X()-donnees.getLC_X())*(1.0/3.0))+donnees.getLC_X())
                {
                    nbPointConsecutif++;
                }
                else
                {
                    nbPointConsecutif = 0;
                }

                if(nbPointConsecutif >= donnees.getTest_K7())
                {
                    donnees.state_X[i].noErreurs.push_back(7);
                    ajouterFailedTests(donnees,7);
                }
            }

            for(int i=0;i < donnees.getDonneesIndividuelle().size();++i)
            {
                if(donnees.getDonneesIndividuelle()[i] < donnees.getLC_X() &&
                        donnees.getDonneesIndividuelle()[i] > ((donnees.getLCI_X()-donnees.getLC_X())*(1.0/3.0))+donnees.getLC_X())
                {
                    nbPointConsecutif++;
                }
                else
                {
                    nbPointConsecutif = 0;
                }

                if(nbPointConsecutif >= donnees.getTest_K7())
                {
                    donnees.state_X[i].noErreurs.push_back(7);
                    ajouterFailedTests(donnees,7);
                }
            }
        }

        //Pour le graphique moving range
        if(donnees.sousGroupe_IMR.size() > donnees.getTest_K7())
        {
            int nbPointConsecutif = 0;

            for(int i=0; i < donnees.sousGroupe_IMR.size();++i)
            {
                if(donnees.sousGroupe_IMR[i].getRange() > donnees.getLC_MR() &&
                        donnees.sousGroupe_IMR[i].getRange() < ((donnees.getLCS_MR()-donnees.getLC_MR())*(1.0/3.0))+donnees.getLC_MR())
                {
                    nbPointConsecutif++;
                }
                else
                {
                    nbPointConsecutif = 0;
                }

                if(nbPointConsecutif >= donnees.getTest_K7())
                {
                    donnees.state_MR[i].noErreurs.push_back(7);
                    //ajouterFailedTests(donnees,7);
                }
            }

            for(int i=0; i < donnees.sousGroupe_IMR.size();++i)
            {
                if(donnees.sousGroupe_IMR[i].getRange() < donnees.getLC_MR() &&
                        donnees.sousGroupe_IMR[i].getRange() > ((donnees.getLCI_MR()-donnees.getLC_MR())*(1.0/3.0))+donnees.getLC_MR())
                {
                    nbPointConsecutif++;
                }
                else
                {
                    nbPointConsecutif = 0;
                }

                if(nbPointConsecutif >= donnees.getTest_K7())
                {
                    donnees.state_MR[i].noErreurs.push_back(7);
                    //ajouterFailedTests(donnees,7);
                }
            }
        }
    }

    //8- huit points consécutifs de la même zone alternant autour de la moyenne (résolution mesure)  --8
    if(donnees.getState_test_K8())
    {
        //Pour le graphique individual
        if(donnees.getDonneesIndividuelle().size() > donnees.getTest_K8())
        {
            int nbPointConsecutif = 0;

            for(int i=0; i < donnees.getDonneesIndividuelle().size(); ++i)
            {
                if(donnees.getDonneesIndividuelle()[i] < (donnees.getLCS_X()-donnees.getLC_X())*(1.0/3.0)+donnees.getLC_X() &&
                        donnees.getDonneesIndividuelle()[i] > (donnees.getLCI_X()-donnees.getLC_X())*(1.0/3.0)+donnees.getLC_X())
                {
                    nbPointConsecutif++;
                }
                else
                {
                    nbPointConsecutif = 0;
                }

                if(nbPointConsecutif >= donnees.getTest_K8())
                {
                    donnees.state_X[i].noErreurs.push_back(8);
                    ajouterFailedTests(donnees,8);
                }
            }
        }

        //Pou le graphique moving range
        if(donnees.sousGroupe_IMR.size() > donnees.getTest_K8())
        {
            int nbPointConsecutif = 0;

            for(int i = 0; i < donnees.sousGroupe_IMR.size(); ++i)
            {
                if(donnees.sousGroupe_IMR[i].getRange() < ((donnees.getLCS_MR()-donnees.getLC_MR())*(1.0/3.0))+donnees.getLC_MR() &&
                        donnees.sousGroupe_IMR[i].getRange() > ((donnees.getLCI_MR()-donnees.getLC_MR())*(1.0/3.0))+donnees.getLC_MR())
                {
                    nbPointConsecutif++;
                }
                else
                {
                    nbPointConsecutif = 0;
                }

                if(nbPointConsecutif >= donnees.getTest_K8())
                {
                    donnees.state_MR[i].noErreurs.push_back(8);
                    //ajouterFailedTests(donnees,8);
                }
            }
        }
    }
}

void AnalyseurCSP::getNormalInfo(DonneesAnalyse &donnees)
{
    //http://boiteaoutils.safran/spip.php?page=methode&lang=fr&id_rubrique=49
    //http://en.wikipedia.org/wiki/Category:Normality_tests
    //Si -1 -> Non-Normal
    //Si  0 -> Passable
    //si  1 -> Normal

    if(isNormalPro(donnees))
    {
        donnees.setIsNormal(1);
    }
    else
    {
        if(  (SkewnessTest(donnees) != 0 && KurtosisTest(donnees) == 0)  || (SkewnessTest(donnees) == 0 && KurtosisTest(donnees) != 0)  )
        {
            if(AndersonDarlingTest(donnees))
            {
                donnees.setIsNormal(0);
            }
            else
            {
                donnees.setIsNormal(-1);
            }
        }
        else
        {
            donnees.setIsNormal(-1);
        }

        if(  SkewnessTest(donnees) == 0 && KurtosisTest(donnees) == 0  )
        {
            if(AndersonDarlingTest(donnees))
            {
                donnees.setIsNormal(1);
            }
            else
            {
                donnees.setIsNormal(-1);
            }
        }


    }

}

//Test de Anderson-Darling
bool AnalyseurCSP::AndersonDarlingTest(DonneesAnalyse &donnees)
{
    //http://www.variation.com/da/help/hs140.htm

    if( (donnees.getPvalue() > 0.05) && (donnees.getAD() < .75) )
    {
        return true;
    }
    else
    {
        return false;
    }

}

//http://www.spcforexcel.com/anderson-darling-test-for-normality
double AnalyseurCSP::AndersonDarlingPValue(int n, double z)
{
    // Returns the AndersonDarling p-value given n and the value of the statistic
    // http://www.minitab.com/fr-CA/support/answers/answer.aspx?id=897&langType=1036

    double ast = z*(1 + 0.75/n + 2.25/(n*n));

    if( 0.600 <= ast && ast < 13 )
    {
        return exp(1.2937 - 5.709*ast + 0.0186*ast*ast);
    }
    if( 0.340 <= ast && ast < 0.600 )
    {
        return exp(0.9177 - 4.279*ast - 1.38*ast*ast);
    }
    if( 0.200 <= ast && ast < 0.340 )
    {
        return 1 - exp(-8.318 + 42.796*ast - 59.938*ast*ast);
    }
    if( ast < 0.200 )
    {
        return 1 - exp(-13.436 + 101.14*ast - 223.73*ast*ast);
    } else
    {
        return -1;
    }


}
double AnalyseurCSP::AndersonDarlingStatistic(DonneesAnalyse &donnees)
{
    //http://google-perftools.googlecode.com/svn-history/r97/trunk/src/tests/sampler_test.cc
    //Ce code n'est pas exact, mais a été quand même utile un peu
    vector<double> donneesIndivitual2(donnees.getDonneesIndividuelle().size());

    copy(donnees.getDonneesIndividuelle().begin(), donnees.getDonneesIndividuelle().end(), donneesIndivitual2.begin());

    std::sort(donneesIndivitual2.begin(), donneesIndivitual2.end());

    double n = donnees.getDonneesIndividuelle().size();
    double m = donnees.getMoyene();
    double s = donnees.getEcartType();

    if(m != 0 && s > 0){
        double ad_sum = 0;
        //on exclu les valeurs des extrémités pour ne pas avoir de CDF de 1...
        for (int i = 1; i < n -1 ; i++)
        {
            double cdf1 = CDF(donneesIndivitual2[i], m, s);
            double cdf2 = CDF(donneesIndivitual2[n-1-i], m, s);
            ad_sum += (2*(i+1) - 1) * (log(cdf1) + log(1 - cdf2));
        }

        double ad_statistic = - n - ad_sum/n;

        return ad_statistic;
    }else
    {
        return -1;
    }
}

int AnalyseurCSP::SkewnessTest(DonneesAnalyse &donnees)
{
    if(donnees.getSkewness() < -0.5)
    {
        return -1;
    }
    else if( (-0.5 <= donnees.getSkewness()) && (donnees.getSkewness() <= 0.5) )
    {
        return 0;
    }
    else if(0.5 < donnees.getSkewness())
    {
        return 1;
    }
    else
    {
        return NULL;
    }
}

int AnalyseurCSP::KurtosisTest(DonneesAnalyse &donnees)
{
    if(donnees.getKurtosis() < -3.5)
    {
        return -1;
    }
    else if( (-3.5 <= donnees.getKurtosis()) && (donnees.getKurtosis() <= 3.5) )
    {
        return 0;
    }
    else if(3.5 < donnees.getKurtosis())
    {
        return 1;
    }
    else
    {
        return NULL;
    }
}

//Test de normale avec un Skewness et le Kurtosis
bool AnalyseurCSP::SkewnessKurtosisAll(DonneesAnalyse &donnees)
{
    //http://www.variation.com/da/help/hs133.htm
    //Test de normale avec un Skewness et le Kurtosis dites:
    //une normale est symetrique alors : -0.5 < Skewness < 0.5
    // et une normale est pas pointue ou aplatie : -3.5 < Kurtosis < 3.5
    if( ( (-0.5 < donnees.getSkewness()) && (donnees.getSkewness() < 0.5) ) && ( (-3.5 < donnees.getKurtosis()) && (donnees.getKurtosis() < 3.5) ) )
    {
        return true;
    }
    else
    {
        return false;
    }

}

//Ce test verifie si les donnees suive une loi normale d'apres les proprietes de la loi normale
bool AnalyseurCSP::isNormalPro(DonneesAnalyse &donnees)
{
    //http://www.ehow.com/how_5284154_test-normality-bell-curve-distribution.html
    //pour faire ca, on verifie si 68.26% des donnees sont entre 1 ecart type de la moyenne,
    //on verifie si 95.44% des donnees sont entre 2 ecart type de la moyenne et
    //on verifie si 99.74% des donnees sont entre 3 ecart type de la moyenne.
    //Proprietes de la loi normale
    size_t size = donnees.getDonneesIndividuelle().size();
    double moyenne = donnees.getMoyene();
    double ecartType = donnees.getEcartType();

    double nb1ET=0;
    double nb2ET=0;
    double nb3ET=0;

    for(int i = 0; i < size; i++)
    {
        if(donnees.getDonneesIndividuelle()[i] >= (moyenne - ecartType) && donnees.getDonneesIndividuelle()[i] <= (moyenne + ecartType))
        {
            nb1ET++;
            nb2ET++;
            nb3ET++;
        }
        else if(donnees.getDonneesIndividuelle()[i] >= (moyenne - 2*ecartType) && donnees.getDonneesIndividuelle()[i] <= (moyenne + 2*ecartType))
        {
            nb2ET++;
            nb3ET++;
        }
        else if(donnees.getDonneesIndividuelle()[i] >= (moyenne - 3*ecartType) && donnees.getDonneesIndividuelle()[i] <= (moyenne + 3*ecartType))
        {
            nb3ET++;
        }

    }

    if(  (0.6800 < nb1ET/size && nb1ET/size < 0.6852) && (0.9500 < nb2ET/size && nb2ET/size < 0.9588) && (0.9900 < nb3ET/size && nb3ET/size < 0.9974)  )
    {
        return true;
    }
    else
    {
        return false;
    }


}

double AnalyseurCSP::KolSmiTest(DonneesAnalyse &donnees)
{
    //ont copie les donnees pour etre sur de ne pas deranger d'autre parti du code
    vector<double> donneesIndivitual2(donnees.getDonneesIndividuelle().size());

    copy(donnees.getDonneesIndividuelle().begin(), donnees.getDonneesIndividuelle().end(), donneesIndivitual2.begin());

    std::sort(donneesIndivitual2.begin(), donneesIndivitual2.end());

    size_t size = donneesIndivitual2.size();

    //Les donnees qui seront generer
    vector<double> donneesTest;

    for(int i = 0; i < size; i++)
    {
        donneesTest.push_back(randn_trig(donnees.getMoyene(), donnees.getEcartType()));
    }

    std::sort(donneesTest.begin(), donneesTest.end());


    return KolmogorovTest(int(size), donneesIndivitual2, int(size), donneesTest);
}

//http://www.dreamincode.net/code/snippet1446.htm
double AnalyseurCSP::randn_notrig(double mu, double sigma)
{
    bool deviateAvailable=false;        //        flag
    static float storedDeviate;                        //        deviate from previous calculation
    double polar, rsquared, var1, var2;

    //        If no deviate has been stored, the polar Box-Muller transformation is
    //        performed, producing two independent normally-distributed random
    //        deviates.  One is stored for the next round, and one is returned.
    if (!deviateAvailable) {

        //        choose pairs of uniformly distributed deviates, discarding those
        //        that don't fall within the unit circle
        do {
            var1=2.0*( double(rand())/double(RAND_MAX) ) - 1.0;
            var2=2.0*( double(rand())/double(RAND_MAX) ) - 1.0;
            rsquared=var1*var1+var2*var2;
        } while ( rsquared>=1.0 || rsquared == 0.0);

        //        calculate polar tranformation for each deviate
        polar=sqrt(-2.0*log(rsquared)/rsquared);

        //        store first deviate and set flag
        storedDeviate=var1*polar;
        deviateAvailable=true;

        //        return second deviate
        return var2*polar*sigma + mu;
    }

    //        If a deviate is available from a previous call to this function, it is
    //        returned, and the flag is set to false.
    else {
        deviateAvailable=false;
        return storedDeviate*sigma + mu;
    }


}
//http://www.dreamincode.net/code/snippet1446.htm
double AnalyseurCSP::randn_trig(double mu, double sigma)
{
    static bool deviateAvailable=false;        //        flag
    static float storedDeviate;                        //        deviate from previous calculation
    double dist, angle;

    //        If no deviate has been stored, the standard Box-Muller transformation is
    //        performed, producing two independent normally-distributed random
    //        deviates.  One is stored for the next round, and one is returned.
    if (!deviateAvailable) {

        //        choose a pair of uniformly distributed deviates, one for the
        //        distance and one for the angle, and perform transformations
        dist=sqrt( -2.0 * log(double(rand()) / double(RAND_MAX)) );
        angle=2.0 * PI * (double(rand()) / double(RAND_MAX));

        //        calculate and store first deviate and set flag
        storedDeviate=dist*cos(angle);
        deviateAvailable=true;

        //        calcaulate return second deviate
        return dist * sin(angle) * sigma + mu;
    }

    //        If a deviate is available from a previous call to this function, it is
    //        returned, and the flag is set to false.
    else {
        deviateAvailable=false;
        return storedDeviate*sigma + mu;
    }

}

//Compare deux vector pour donner un pourcentage de normality (?)
double AnalyseurCSP::KolmogorovTest(int na, vector<double> a, int nb, vector<double> b)
{
    //  Statistical test whether two one-dimensional sets of points are compatible
    //  with coming from the same parent distribution, using the Kolmogorov test.
    //  That is, it is used to compare two experimental distributions of unbinned data.
    //
    //  Input:
    //  a,b: One-dimensional arrays of length na, nb, respectively.
    //       The elements of a and b must be given in ascending order.
    //  option is a character string to specify options
    //         "D" Put out a line of "Debug" printout
    //         "M" Return the Maximum Kolmogorov distance instead of prob
    //
    //  Output:
    // The returned value prob is a calculated confidence level which gives a
    // statistical test for compatibility of a and b.
    // Values of prob close to zero are taken as indicating a small probability
    // of compatibility. For two point sets drawn randomly from the same parent
    // distribution, the value of prob should be uniformly distributed between
    // zero and one.
    //   in case of error the function return -1
    //   If the 2 sets have a different number of points, the minimum of
    //   the two sets is used.
    //
    // Method:
    // The Kolmogorov test is used. The test statistic is the maximum deviation
    // between the two integrated distribution functions, multiplied by the
    // normalizing factor (rdmax*sqrt(na*nb/(na+nb)).


    double prob = -1;
    //     Constants needed
    double rna = na;
    double rnb = nb;
    double sa  = 1./rna;
    double sb  = 1./rnb;
    double rdiff = 0;
    double rdmax = 0;
    int ia = 0;
    int ib = 0;

    //    Main loop over point sets to find max distance
    //    rdiff is the running difference, and rdmax the max.
    bool ok = false;

    for (int i = 0; i < na+nb; i++) {
        if (a[ia] < b[ib]) {
            rdiff -= sa;
            ia++;
            if (ia >= na)
            {
                ok = true;
                break;
            }
        } else if (a[ia] > b[ib])
        {
            rdiff += sb;
            ib++;
            if (ib >= nb)
            {
                ok = true;
                break;
            }
        } else
        {
            // special cases for the ties
            double x = a[ia];
            while(a[ia] == x && ia < na -1)
            {
                rdiff -= sa;
                ia++;
            }
            while(b[ib] == x && ib < nb-1)
            {
                rdiff += sb;
                ib++;
            }
            if (ia >= na-1)
            {
                ok = true;
                break;
            }
            if (ib >= nb-1)
            {
                ok = true;
                break;
            }
        }
        rdmax = std::max(rdmax,abs(rdiff));
    }

    if (ok) {
        rdmax = std::max(rdmax,abs(rdiff));
        double z = rdmax * sqrt(rna*rnb/(rna+rnb));
        prob = KolmogorovProb(z);
        return prob;
    }else
    {
        return -1;
    }

}

double AnalyseurCSP::KolmogorovProb(double z)
{
    double fj[4] = {-2,-8,-18,-32}, r[4];
    const double w = 2.50662827;
    // c1 - -pi**2/8, c2 = 9*c1, c3 = 25*c1
    const double c1 = -1.2337005501361697;
    const double c2 = -11.103304951225528;
    const double c3 = -30.842513753404244;

    double u = abs(z);
    double p;
    if (u < 0.2) {
        p = 1;
    } else if (u < 0.755) {
        double v = 1./(u*u);
        p = 1 - w*(exp(c1*v) + exp(c2*v) + exp(c3*v))/u;
    } else if (u < 6.8116) {
        r[1] = 0;
        r[2] = 0;
        r[3] = 0;
        double v = u*u;
        int maxj = std::max(1,Nint(3./u));
        for (int j = 0; j < maxj; j++) {
            r[j] = exp(fj[j]*v);
        }
        p = 2*(r[0] - r[1] +r[2] - r[3]);
    } else {
        p = 0;
    }
    return p;

}

int AnalyseurCSP::Nint(float x)
{
    // Round to nearest integer. Rounds half integers to the nearest
    // even integer.

    int i;
    if (x >= 0) {
        i = int(x + 0.5);
        if (x + 0.5 == float(i) && i & 1) i--;
    } else {
        i = int(x - 0.5);
        if (x - 0.5 == float(i) && i & 1) i++;

    }
    return i;

}

int AnalyseurCSP::Nint(double x)
{
    // Round to nearest integer. Rounds half integers to the nearest
    // even integer.

    int i;
    if (x >= 0) {
        i = int(x + 0.5);
        if (x + 0.5 == double(i) && i & 1) i--;
    } else {
        i = int(x - 0.5);
        if (x - 0.5 == double(i) && i & 1) i++;

    }
    return i;

}

double AnalyseurCSP::CDF(double x, double m, double s)
{
    if (s != 0)
    {
        boost::math::normal myNormal(m, s);
        return cdf(myNormal, x);
    }

    return 0;
}

void AnalyseurCSP::BoxCoxTransformation(DonneesAnalyse &donnees, vector<DonneesAnalyse*> &output)
{

    double l = -5.0;
    int i = 0;
    while(l <= 5.0)
    {
        DonneesAnalyse * donneesSB = new DonneesAnalyse();
        donneesSB->setPopOuEchan(false);
        for(int n = 0; n < donnees.getDonneesIndividuelle().size(); n++)
        {
            if(l == 0)
            {
                //donneesSB->getDonneesIndividuelle().push_back(log10( donnees.getDonneesIndividuelle()[n]));
            }
            else
            {
                //donneesSB->getDonneesIndividuelle().push_back(pow( donnees.getDonneesIndividuelle()[n] , l));
            }
        }

        donneesSB->setLoTol(donnees.getLoTol());
        donneesSB->setUpTol(donnees.getUpTol());
        donneesSB->setTolNominal(donnees.getTolNominal());
        donneesSB->setLimiteInf(donnees.getLimiteInf());
        donneesSB->setLimiteSup(donnees.getLimiteSup());
        //calculIMR((*donneesSB));
        calculVariance((*donneesSB));
        calculEcartType((*donneesSB));

        output.push_back(donneesSB);

        l=l+1;
        i++;
    }
}

void AnalyseurCSP::ajouterFailedTests(DonneesAnalyse &donnees, int noErreur)
{
    QVector<int> failedTest = donnees.getFailedTests();

    //ajoute l'erreur si pas déjà présente et trie le tableau
    if (!failedTest.contains(noErreur))
    {
        failedTest.append(noErreur);
    }
    qSort(failedTest);

    donnees.setFailedTests(failedTest);
}
