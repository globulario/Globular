#include "SousGroupe.h"
#if defined(WIN32) && defined(_DEBUG)
#define DEBUG_NEW new( _NORMAL_BLOCK, __FILE__, __LINE__ )
#define new DEBUG_NEW
#endif
//Constructeur et destructeur
SousGroupe::SousGroupe() :
    range(0.0f),
    moyenne(0.0f),
    state(false)
{

}

SousGroupe::~SousGroupe()
{}

//Accesseurs
double SousGroupe::getRange()
{
    return range;
}
double SousGroupe::getMoyenne()
{
    return moyenne;
}
bool SousGroupe::getState()
{
    return state;
}

//Mutateurs
void SousGroupe::setRange(double range)
{
    this->range = range;
}
void SousGroupe::setMoyenne(double moyenne)
{
    this->moyenne = moyenne;
}
void SousGroupe::setState(bool state)
{
    this->state = state;
}

void SousGroupe::read(const QJsonObject &json){
    /** nothing to do here **/
}

void SousGroupe::write(QJsonObject &json) const {

    json["range"] = this->range;
    json["moyenne"] = this->moyenne;
    json["state"] = this->state;

    // Now the values.
    QJsonArray donnees;

    for(size_t i=0; i<this->donnees.size(); i++){
        double donnee = this->donnees[i];
        donnees.append(donnee);
    }

    // Set the array of values.
    json["donnees"] = donnees;
}
