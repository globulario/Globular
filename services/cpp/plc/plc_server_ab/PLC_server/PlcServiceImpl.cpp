#include "PlcServiceImpl.h"
#include <string>
#include <fstream>
#include <iostream>
#include <sstream>
#include <map>
#include <cstddef>
#include <bitset>         // std::bitset
#include <math.h>       /* ceil */

//  https://github.com/nlohmann/json
#include "json.hpp"

#ifdef WIN32
#include <windows.h>
std::string getexepath()
{
	char result[MAX_PATH];
	return std::string(result, GetModuleFileName(NULL, result, MAX_PATH));
}

void sleep(unsigned milliseconds)
{
	Sleep(milliseconds);
}

#else
#include <limits.h>
#include <unistd.h>
#include <linux/limits.h>

std::string getexepath()
{
	char result[PATH_MAX];
	ssize_t count = readlink("/proc/self/exe", result, PATH_MAX);
	return std::string(result, (count > 0) ? count : 0);
}

#endif // WIN32

// for convenience
//using json = nlohmann::json;
using std::cout;

// The constructor
PlcServiceImpl::PlcServiceImpl(unsigned int defaultPort, unsigned int defaultProxy) :
	defaultPort(defaultPort),
	defaultProxy(defaultProxy),
	name("plc_server_ab"),
	domain("localhost"),
	publisher_id("localhost"),
	allow_all_origins(true),
	allowed_origins(""),
	tls(false)
{

	// first of all I will try to open the configuration from the file.
	std::ifstream inFile;
	std::string execPath = getexepath();
#ifdef WIN32
	std::size_t lastIndex = execPath.find_last_of("/\\");
	std::string configPath = execPath.substr(0, lastIndex) + "\\config.json";
#else
	std::size_t lastIndex = execPath.find_last_of("/");
	std::string configPath = execPath.substr(0, lastIndex) + "/config.json";
#endif
	inFile.open(configPath); //open the input file

	if (inFile.good()) {

		std::stringstream strStream;
		strStream << inFile.rdbuf(); //read the file
		std::string jsonStr = strStream.str(); //str holds the content of the file

		// Parse the json file.
		auto j = nlohmann::json::parse(jsonStr);

		// Now I will initialyse the value from the configuration file.
		this->publisher_id = j["PublisherId"];
		this->version = j["Version"];
		this->keep_up_to_date = j["KeepUpToDate"];
		this->allow_all_origins = j["AllowAllOrigins"];
		this->cert_authority_trust = j["CertAuthorityTrust"];
		this->keep_alive = j["KeepAlive"];
		this->cert_file = j["CertFile"];
		this->domain = j["Domain"];
		this->key_file = j["KeyFile"];
		this->name = j["Name"];
		this->defaultPort = j["Port"];
		this->defaultProxy = j["Proxy"];
		this->tls = j["TLS"];
		// can be a list of string
		this->allowed_origins = j["AllowedOrigins"];

		// Now the connections.
		auto connections = j["Connections"];
		for (auto connection : connections) {
			// Here i will initialyse the connection.

			Connection c;
			c.id = connection["Id"];
			c.ip = connection["Ip"];
			c.slot = connection["Slot"];
			c.timeout = connection["Timeout"];
			c.protocol = ProtocolType(connection["Protocol"]);
			c.cpu = CpuType(connection["Cpu"]);
			c.portType = PortType(connection["PortType"]);

			this->connections[c.id] = c;
		}
	}
	else {
		// save the default configuration.
		this->save();
	}

}

// The desctructor.
PlcServiceImpl::~PlcServiceImpl() {

}

void PlcServiceImpl::save() {
	/**/
	nlohmann::json j;
	j["PublisherId"] = this->publisher_id;
	j["Version"] = this->version;
	j["KeepUpToDate"] = this->keep_up_to_date;
	j["KeepAlive"] = this->keep_alive;
	j["AllowAllOrigins"] = this->allow_all_origins;
	j["AllowedOrigins"] = this->allowed_origins; // empty string
	j["CertAuthorityTrust"] = this->cert_authority_trust;
	j["CertFile"] = this->cert_file;
	j["Domain"] = this->domain;
	j["KeyFile"] = this->key_file;
	j["Name"] = this->name;
	j["Port"] = this->defaultPort;
	j["Name"] = this->name;
	j["Protocol"] = "grpc";
	j["Proxy"] = this->defaultProxy;
	j["TLS"] = this->tls;

	// Create connection objects.
	auto connections = nlohmann::json::array();
	for (auto& kv : this->connections) {
		auto c = kv.second;
		auto connection = nlohmann::json::object();
		// Here i will initialyse the connection.
		connection["id"] = c.id;
		connection["ip"] = c.ip;
		connection["slot"] = c.slot;
		connection["timeout"] = c.timeout;
		connection["protocol"] = c.protocol;
		connection["cpu"] = c.cpu;
		connection["portType"] = c.portType;
		if (c.save) {
			connections.push_back(connection);
		}
	}

	// Set the connections...
	j["Connections"] = connections;
	std::string execPath = getexepath();
#ifdef WIN32
	std::size_t lastIndex = execPath.find_last_of("/\\");
	std::string configPath = execPath.substr(0, lastIndex) + "\\config.json";
#else
	std::size_t lastIndex = execPath.find_last_of("/");
	std::string configPath = execPath.substr(0, lastIndex) + "/config.json";
#endif
	std::ofstream file;
	file.open(configPath);
	file << j.dump();
	file.close();
}

// Create a connection.
::grpc::Status PlcServiceImpl::CreateConnection(::grpc::ServerContext* context, const ::plc::CreateConnectionRqst* request, ::plc::CreateConnectionRsp* response) {

	::grpc::Status s;

	// Here I will get the connection object from the request.
	auto connection = request->connection();

	Connection c;
	c.id = connection.id();
	c.ip = connection.ip();
	c.slot = connection.slot();
	c.timeout = connection.timeout();
	c.portType = PortType(connection.protocol());
	c.cpu = CpuType(connection.cpu());
	c.protocol = ProtocolType(connection.protocol());
	c.save = request->save();

	// set the connection in the map.
	this->connections[connection.id()] = c;

	// save the config with the connection in it.
	if (request->save()) {
		this->save();
	}

	return ::grpc::Status::OK;
}

std::string PlcServiceImpl::GetProtocol(Connection connection) {
	//cout << "protocol " << connection.protocol << std::endl;
	// Get the connection id.
	if (connection.protocol == ProtocolType::AB_EIP) {
		return "ab_eip";
	}
	if (connection.protocol == ProtocolType::AB_CIP) {
		return "ab_cip";
	}

	return "";
}

std::string PlcServiceImpl::GetPortType(Connection connection) {
	//cout << "portType " << connection.portType << std::endl;

	// Get the connection id.
	if (connection.portType == PortType::BACKPLANE) {
		return "1";
	}
	if (connection.portType == PortType::NET_ETHERNET) {
		return "2";
	}
	if (connection.portType == PortType::DH_PLUS_CHANNEL_A) {
		return "2";
	}
	if (connection.portType == PortType::DH_PLUS_CHANNEL_B) {
		return "2";
	}
	if (connection.portType == PortType::SERIAL) {
		return "3";
	}

	return "";
}

std::string PlcServiceImpl::GetCpuType(Connection connection) {
	//cout << "cpuType " <<  connection.cpu << std::endl;

	if (connection.cpu == CpuType::PLC) {
		return "PLC";
	}

	if (connection.cpu == CpuType::PLC5) {
		return "PLC5";
	}

	if (connection.cpu == CpuType::SLC) {
		return "SLC";
	}

	if (connection.cpu == CpuType::SLC500) {
		return "SLC500";
	}

	if (connection.cpu == CpuType::MICROLOGIX) {
		return "MICROLOGIX";
	}

	if (connection.cpu == CpuType::MLGX) {
		return "MLGX";
	}

	if (connection.cpu == CpuType::COMPACTLOGIX) {
		return "COMPACTLOGIX";
	}

	if (connection.cpu == CpuType::CLGX) {
		return "CLGX";
	}

	if (connection.cpu == CpuType::LGX) {
		return "LGX";
	}

	if (connection.cpu == CpuType::CONTROLLOGIX) {
		return "CONTROLLOGIX";
	}

	if (connection.cpu == CpuType::CONTROLOGIX) {
		return "CONTROLOGIX";
	}

	if (connection.cpu == CpuType::FLEXLOGIX) {
		return "FLEXLOGIX";
	}

	if (connection.cpu == CpuType::FLGX) {
		return "FLGX";
	}

	return "";
}

int PlcServiceImpl::GetTypeSize(TagType type) {
	if (type == TagType::SINT_TAG_TYPE) {
		return 1;
	}

	if (type == TagType::INT_TAG_TYPE) {
		return 2;
	}

	if (type == TagType::BOOL_TAG_TYPE || type == TagType::DINT_TAG_TYPE || type == TagType::REAL_TAG_TYPE) {
		return 4;
	}

	if (type == TagType::LINT_TAG_TYPE || type == TagType::LREAL_TAG_TYPE) {
		return 8;
	}

	return 0;
}

// Get Communication path in a format compatible with the CIP library (only until cpu parameter)
// TAG_PATH= "protocol=ab_eip&gateway=192.168.1.207&path=1,0&cpu=LGX&elem_count=1&elem_size=4&name=Real"
std::string PlcServiceImpl::GetCpuPath(Connection connection)
{
	return "protocol=" + this->GetProtocol(connection) + "&gateway=" + connection.ip + "&path=" + this->GetPortType(connection) + "," + std::to_string(connection.slot) + "&cpu=" + this->GetCpuType(connection);
};


::grpc::Status PlcServiceImpl::GetConnection(::grpc::ServerContext* context, const ::plc::GetConnectionRqst* request, ::plc::GetConnectionRsp* response) {
	std::map<std::string, Connection>::iterator it = this->connections.find(request->id());
	if (it == this->connections.end()) {
		return ::grpc::Status(::grpc::StatusCode::INTERNAL, grpc::string("ERROR: Connection not found!"));
	}

	response->mutable_connection()->set_id((*it).second.id);
	response->mutable_connection()->set_ip((*it).second.ip);
	response->mutable_connection()->set_cpu(::plc::CpuType((*it).second.cpu));
	response->mutable_connection()->set_porttype(::plc::PortType((*it).second.portType));
	response->mutable_connection()->set_protocol(::plc::ProtocolType((*it).second.protocol));
	response->mutable_connection()->set_slot((*it).second.slot);
	response->mutable_connection()->set_timeout((*it).second.timeout);

	return ::grpc::Status::OK;
}

// Close the connection and release all open tags. *
::grpc::Status PlcServiceImpl::CloseConnection(::grpc::ServerContext* context, const ::plc::CloseConnectionRqst* request, ::plc::CloseConnectionRsp* response) {
	std::map<std::string, Connection>::iterator it = this->connections.find(request->connection_id());
	if (it == this->connections.end()) {
		return ::grpc::Status(::grpc::StatusCode::INTERNAL, grpc::string("Connection " + request->connection_id() + " not exist!"));
	}

	// remove all tags.
	for (auto it = this->tags.begin(); it != this->tags.end(); ++it) {
		plc_tag_destroy(it->second);
		this->tags.erase(it);
	}

	return ::grpc::Status::OK;
}

// Delete a connection.
::grpc::Status PlcServiceImpl::DeleteConnection(::grpc::ServerContext* context, const ::plc::DeleteConnectionRqst* request, ::plc::DeleteConnectionRsp* response) {
	std::map<std::string, Connection>::iterator it = this->connections.find(request->id());
	if (it != this->connections.end()) {
		this->connections.erase(it);
		this->save();
	}
	return ::grpc::Status::OK;
}

// Add tag in the list and create pointer to the plc tag
// It assume the SetCommParam method was called before to set all communication parameters
// return: index of the tag in the list, or -1 if had a problem
::grpc::Status PlcServiceImpl::OpenTag(std::string connectionId, std::string name, TagType type, int elementCount) {
	std::map<std::string, Connection>::iterator connectionIt = this->connections.find(connectionId);
	if (connectionIt == this->connections.end()) {
		return ::grpc::Status(::grpc::StatusCode::INTERNAL, grpc::string("Connection " + connectionId + " not exist!"));
	}


	// first of all I will get the tag from the map of tags.
	std::map<std::string, int>::const_iterator tagIt = this->tags.find(connectionId + ":" + name);

	// Create the path.
	std::string path = this->GetCpuPath(connectionIt->second);
	std::string tagPath = path + "&elem_count=" + std::to_string(elementCount) + "&elem_size=" + std::to_string(this->GetTypeSize(type)) + "&name=" + name;

	// test if the path has change if it does I need to re-open the tag before reading it.
	std::map<std::string, std::string>::const_iterator pathIt = this->paths.find(connectionId + ":" + name);
	bool hasChange = false;
	if (pathIt != this->paths.end()) {
		hasChange = pathIt->second != tagPath;
	}

	if (tagIt == this->tags.end() || hasChange) {

		// Keep the path.
		this->paths[connectionId + ":" + name] = tagPath;

		// Create tag in CIP and check for status
		int tag = plc_tag_create(tagPath.c_str(), int(connectionIt->second.timeout));

		if (!tag) {
			return ::grpc::Status(::grpc::StatusCode::INTERNAL, grpc::string("ERROR: Could not create tag"));
		}
		//std::cout << "tag path : " << tagPath << std::endl;

		// let the connect succeed
		while (plc_tag_status(tag) == PLCTAG_STATUS_PENDING) {
			sleep(5);
		}

		if (plc_tag_status(tag) != PLCTAG_STATUS_OK) {
			plc_tag_destroy(tag);
			return ::grpc::Status(::grpc::StatusCode::INTERNAL, grpc::string("Error setting up tag internal state"));
		}

		// Save the tag to furter use.
		this->tags[connectionId + ":" + name] = tag;
	}

	return ::grpc::Status::OK;
}

// * Get tags *
::grpc::Status PlcServiceImpl::ReadTag(::grpc::ServerContext* context, const ::plc::ReadTagRqst* request, ::plc::ReadTagRsp* response) {

	std::map<std::string, Connection>::iterator connectionIt = this->connections.find(request->connection_id());
	if (connectionIt == this->connections.end()) {
		return ::grpc::Status(::grpc::StatusCode::INTERNAL, grpc::string("Connection " + request->connection_id() + " not exist!"));
	}

	// Now I will read the tag value.
	TagType type = TagType(request->type());
	int offset = request->offset();
	int length = request->length();
	int elementSize = GetTypeSize(type);

	// The size of element to open.
	int size = offset + length;

	if (type == TagType::BOOL_TAG_TYPE) {
		size = std::ceil((float)(size / 32.0f)); // bool contain 32 bit.
	}

	auto status = OpenTag(request->connection_id(), request->name(), type, size);

	// If there is an error
	if (!status.ok()) {
		return status;
	}

	// first of all I will get the tag from the map of tags.
	std::map<std::string, int>::const_iterator tagIt = this->tags.find(connectionIt->first + ":" + request->name());
	if (tagIt == this->tags.end()) {
		return ::grpc::Status(::grpc::StatusCode::INTERNAL, grpc::string("Tag with name " + request->name() + " is not define for connection " + connectionIt->first));
	}

	// Ask for a data to the plc.
	int  rc = plc_tag_read(tagIt->second, int(connectionIt->second.timeout)); // read tag into CPI buffer
	if (rc != PLCTAG_STATUS_OK) {

		// remove the tag  from the local map.
		plc_tag_destroy(tagIt->second);
		this->tags.erase(tagIt);

		std::stringstream ss;
		ss << "ERROR: Unable to read the data! Got error code " << plc_tag_decode_error(rc);
		return ::grpc::Status(::grpc::StatusCode::INTERNAL, grpc::string(ss.str()));
	}; // if not status ok - end

	std::string values = "[";

	if (type == TagType::BOOL_TAG_TYPE) {
		for (int i = offset; i < offset + length; i++) {
			auto b = plc_tag_get_bit(tagIt->second, i);
			values += std::to_string(b);
			// json array...
			if (i < length + offset - 1) {
				values += ", ";
			}
		}
	}
	else {
		for (int i = offset; i < offset + length; i++) {
			auto index = i * elementSize;
			if (type == TagType::DINT_TAG_TYPE) {
				if (request->unsigned_()) {
					values += std::to_string(plc_tag_get_uint32(tagIt->second, index));
				}
				else {
					values += std::to_string(plc_tag_get_int32(tagIt->second, index));
				}
			}
			else if (type == TagType::INT_TAG_TYPE) {
				if (request->unsigned_()) {
					values += std::to_string(plc_tag_get_uint16(tagIt->second, index));
				}
				else {
					values += std::to_string(plc_tag_get_int16(tagIt->second, index));
				}
			}
			else if (type == TagType::SINT_TAG_TYPE) {
				if (request->unsigned_()) {
					values += std::to_string(plc_tag_get_uint8(tagIt->second, index));
				}
				else {
					values += std::to_string(plc_tag_get_int8(tagIt->second, index));
				}
			}
			else if (type == TagType::LINT_TAG_TYPE) {
				if (request->unsigned_()) {
					values += std::to_string(plc_tag_get_uint64(tagIt->second, index));
				}
				else {
					values += std::to_string(plc_tag_get_int64(tagIt->second, index));
				}
			}
			else if (type == TagType::REAL_TAG_TYPE) {
				values += std::to_string(plc_tag_get_float32(tagIt->second, index));
			}
			else if (type == TagType::LREAL_TAG_TYPE) {
				values += std::to_string(plc_tag_get_float64(tagIt->second, index));
			}

			// json array...
			if (i < length + offset - 1) {
				values += ", ";
			}
		}
	}

	values += "]";

	// Set the response.
	response->set_values(values);

	return ::grpc::Status::OK;
}

// * Write tag *
::grpc::Status PlcServiceImpl::WriteTag(::grpc::ServerContext* context, const ::plc::WriteTagRqst* request, ::plc::WriteTagRsp* response) {
	std::map<std::string, Connection>::iterator connectionIt = this->connections.find(request->connection_id());
	if (connectionIt == this->connections.end()) {
		return ::grpc::Status(::grpc::StatusCode::INTERNAL, grpc::string("Connection " + request->connection_id() + " not exist!"));
	}

	// Now I will set the actual value.
	TagType type = TagType(request->type());
	int offset = request->offset();
	int length = request->length();
	int elementSize = GetTypeSize(type);
	
	// The size of element to open.
	int size = offset + length;

	if (type == TagType::BOOL_TAG_TYPE) {
		size = std::ceil((float)(size / 32.0f)); // bool contain 32 bit.
	}

	auto status = OpenTag(request->connection_id(), request->name(), type, size);
	// If there is an error
	if (!status.ok()) {
		return status;
	}

	// first of all I will get the tag from the map of tags.
	std::map<std::string, int>::const_iterator tagIt = this->tags.find(connectionIt->first + ":" + request->name());
	if (tagIt == this->tags.end()) {
		return ::grpc::Status(::grpc::StatusCode::INTERNAL, grpc::string("Tag with name " + request->name() + " is not define for connection " + connectionIt->first));
	}

	int rc = plc_tag_write(tagIt->second, int(connectionIt->second.timeout));

	if (rc != PLCTAG_STATUS_OK) {
		plc_tag_destroy(tagIt->second);
		this->tags.erase(tagIt);

		std::stringstream ss;
		ss << "ERROR: Unable to write the data! Got error code " << rc;
		return ::grpc::Status(::grpc::StatusCode::INTERNAL, grpc::string(ss.str()));
	};

	// Here I will use the incomming
	auto  values = nlohmann::json::parse(request->values());

	if (type == TagType::BOOL_TAG_TYPE) {
		for (int i = offset, j=0; i < offset + length; i++, j++) {
			auto& v = values[j]; // get the value to set.
		    plc_tag_set_bit(tagIt->second, i, v.get<uint8_t>());
		}
	}
	else {
		for (int i = offset, j = 0; i < offset + length; i++, j++) {
			auto index = i * elementSize;
			auto& v = values[j];

			if (type == TagType::SINT_TAG_TYPE) {
				if (request->unsigned_()) {
					uint8_t value = v.get<uint8_t>();
					plc_tag_set_uint8(tagIt->second, index, value);
				}
				else {
					int8_t value = v.get<int8_t>();
					plc_tag_set_int8(tagIt->second, index, value);
				}
			}
			else if (type == TagType::INT_TAG_TYPE) {
				if (request->unsigned_()) {
					uint16_t value = v.get<uint16_t>();
					plc_tag_set_uint16(tagIt->second, index, value);
				}
				else {
					int16_t value = v.get<int16_t>();
					plc_tag_set_int16(tagIt->second, index, value);
				}
			}
			else if (type == TagType::DINT_TAG_TYPE) {
				if (request->unsigned_()) {
					uint32_t value = v.get<uint32_t>();
					plc_tag_set_uint32(tagIt->second, index, value);
				}
				else {
					int32_t value = v.get<int32_t>();
					plc_tag_set_int32(tagIt->second, index, value);
				}
			}
			else if (type == TagType::LINT_TAG_TYPE) {
				if (request->unsigned_()) {
					uint64_t value = v.get<uint64_t>();
					plc_tag_set_uint64(tagIt->second, index, value);
				}
				else {
					int64_t value = v.get<int64_t>();
					plc_tag_set_int64(tagIt->second, index, value);
				}
			}
			else if (type == TagType::REAL_TAG_TYPE) {
				float value = v.get<float>();
				//std::cout << "write value [" << i << "] at index " << index << "with value " << value << std::endl;
				plc_tag_set_float32(tagIt->second, index, value);
			}
			else if (type == TagType::LREAL_TAG_TYPE) {
				double value = v.get<double>();
				plc_tag_set_float64(tagIt->second, index, value);
			}
			else {
				return ::grpc::Status(::grpc::StatusCode::INTERNAL, grpc::string("Tag type not found must be one of float, dint, sint or real"));
			}
		}
	}

	response->set_result(true);
	return ::grpc::Status::OK;
}
