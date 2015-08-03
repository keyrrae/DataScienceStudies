exchangeRate <- read.csv("./USDtoRMB.csv",header=TRUE, sep=",")
exchangeRateTS <- ts(exchangeRate$Mid, frequency=365, start=c(2012,1))
plot.ts(exchangeRateTS)

install.pakcages("TTR")
library("TTR")

exchangeRateSMA5 <- SMA(exchangeRateTS, n=5)
plot.ts(exchangeRateSMA5)

exchangeRateSMA10 <- SMA(exchangeRateTS, n=10)
plot.ts(exchangeRateSMA10)

exchangeRateSMA30 <- SMA(exchangeRateTS, n=30)
plot.ts(exchangeRateSMA30)

exchangeRateDecomposed <- decompose(exchangeRateTS)
plot(exchangeRateDecomposed)
