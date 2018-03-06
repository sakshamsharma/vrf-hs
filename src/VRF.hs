{-# LANGUAGE ForeignFunctionInterface #-}

module VRF where

import           Foreign.C.String
import           Foreign.Marshal.Alloc
import           Prelude

foreign import ccall "DoVRF" go_DoVRF :: CString -> IO CString
foreign import ccall "VerifyVRF" go_VerifyVRF :: CString -> IO CString

runStrFxn :: (CString -> IO CString) -> String -> IO String
runStrFxn f input = do
  cinput <- newCString input
  coutput <- f cinput
  res <- peekCString coutput
  _ <- free cinput
  _ <- free coutput
  return res

ffiDoVRF :: String -> IO String
ffiDoVRF = runStrFxn go_DoVRF

ffiVerifyVRF :: String -> IO String
ffiVerifyVRF = runStrFxn go_VerifyVRF
