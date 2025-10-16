// Package main User Service API
//
// This is the user service for the DDD microservices architecture.
// It provides user management functionality including registration, authentication,
// profile management, and role-based access control.
//
//	@title			User Service API
//	@version		1.0
//	@description	User Service API for DDD Microservices
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:8080
//	@BasePath	/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
//
// @tag.name			Authentication
// @tag.description	User authentication endpoints
// @tag.name			User Management
// @tag.description	User profile and management endpoints
// @tag.name			Admin
// @tag.description	Administrative endpoints (admin only)
package main
