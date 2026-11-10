/*
 * Copyright (c) 2024. Abstrium SAS <team (at) pydio.com>
 * This file is part of Pydio Cells.
 *
 * Pydio Cells is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Pydio Cells is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Pydio Cells.  If not, see <http://www.gnu.org/licenses/>.
 *
 * The latest code can be found at <https://pydio.com>.
 */

package errors

import (
	tozd "gitlab.com/tozd/go/errors"
)

var (
	// CellsError is the parent of all errors
	CellsError = tozd.Base("cells")

	StatusInternalServerError = tozd.BaseWrap(CellsError, "internal")
	StatusForbidden           = tozd.BaseWrap(CellsError, "forbidden")
	StatusUnauthorized        = tozd.BaseWrap(CellsError, "unauthorized")
	StatusTooManyRequests     = tozd.BaseWrap(CellsError, "too many requests")
	StatusRequestTimeout      = tozd.BaseWrap(CellsError, "timeout")
	StatusBadRequest          = tozd.BaseWrap(CellsError, "bad request")
	StatusConflict            = tozd.BaseWrap(CellsError, "conflict")
	StatusPreconditionFailed  = tozd.BaseWrap(CellsError, "precondition failed")
	StatusNotImplemented      = tozd.BaseWrap(CellsError, "not implemented")
	StatusServiceUnavailable  = tozd.BaseWrap(CellsError, "service unavailable")
	StatusNotFound            = tozd.BaseWrap(CellsError, "not found")
	StatusCancelled           = tozd.BaseWrap(CellsError, "cancelled")
	StatusAborted             = tozd.BaseWrap(CellsError, "aborted")
	StatusOutOfRange          = tozd.BaseWrap(CellsError, "out of range")
	StatusDataLoss            = tozd.BaseWrap(CellsError, "data loss")
	StatusLocked              = tozd.BaseWrap(CellsError, "locked")
	StatusQuotaReached        = tozd.BaseWrap(CellsError, "quota reached")

	InvalidParameters = RegisterBaseSentinel(StatusBadRequest, "invalid parameters")

	DAO            = RegisterBaseSentinel(StatusInternalServerError, "dao")
	ResolveError   = RegisterBaseSentinel(DAO, "dao resolution failed")
	LogDAOResolver = RegisterBaseSentinel(DAO, "logger storage not ready")
	SqlDAO         = RegisterBaseSentinel(DAO, "sql")
	MongoDAO       = RegisterBaseSentinel(DAO, "mongo")
	BoltDAO        = RegisterBaseSentinel(DAO, "bolt")
	BleveDAO       = RegisterBaseSentinel(DAO, "bleve")

	ServiceError            = RegisterBaseSentinel(StatusInternalServerError, "service error")
	ServiceStartError       = RegisterBaseSentinel(ServiceError, "service start error")
	ServiceNoServerAttached = RegisterBaseSentinel(ServiceError, "no server attached")
	HealthCheckError        = RegisterBaseSentinel(StatusServiceUnavailable, "healthcheck service error")
)
