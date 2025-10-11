/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package logger 

type Logger interface {
  Info(msg string)
  Warn(msg string)
  Error(msg string)
  Debug(msg string)
  SetVerbose(verbose bool)
}
