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
 * PostInitialize200Response
 *
 * @package IsuRide\Model
 * @author  OpenAPI Generator team
 * @link    https://github.com/openapitools/openapi-generator
 */
class PostInitialize200Response extends BaseModel
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
  "required" : [ "language" ],
  "properties" : {
    "language" : {
      "type" : "string",
      "description" : "実装言語\n- go\n- perl\n- php\n- python\n- ruby\n- rust\n- node\n"
    }
  }
}
SCHEMA;
}
