<?php

/**
 * ISURIDE API Specification
 * PHP version 7.4
 *
 * @package IsuRide
 * @author  OpenAPI Generator team
 * @link    https://github.com/openapitools/openapi-generator
 */

/**
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 * The version of the OpenAPI document: 1.0
 * Generated by: https://github.com/openapitools/openapi-generator.git
 */

/**
 * NOTE: This class is auto generated by the openapi generator program.
 * https://github.com/openapitools/openapi-generator
 */
namespace IsuRide\Model;

use IsuRide\BaseModel;

/**
 * AppChairStatsRecentRidesInner
 *
 * @package IsuRide\Model
 * @author  OpenAPI Generator team
 * @link    https://github.com/openapitools/openapi-generator
 */
class AppChairStatsRecentRidesInner extends BaseModel
{
    /**
     * @var string Models namespace.
     * Can be required for data deserialization when model contains referenced schemas.
     */
    protected const MODELS_NAMESPACE = '\IsuRide\Model';

    /**
     * @var string Constant with OAS schema of current class.
     * Should be overwritten by inherited class.
     */
    protected const MODEL_SCHEMA = <<<'SCHEMA'
{
  "required" : [ "destination_coordinate", "distance", "duration", "evaluation", "id", "pickup_coordinate" ],
  "properties" : {
    "id" : {
      "type" : "string",
      "description" : "ライドID"
    },
    "pickup_coordinate" : {
      "$ref" : "#/components/schemas/Coordinate"
    },
    "destination_coordinate" : {
      "$ref" : "#/components/schemas/Coordinate"
    },
    "distance" : {
      "type" : "integer",
      "description" : "移動距離"
    },
    "duration" : {
      "type" : "integer",
      "description" : "移動時間 (ミリ秒)",
      "format" : "int64"
    },
    "evaluation" : {
      "type" : "integer",
      "description" : "評価"
    }
  }
}
SCHEMA;
}
