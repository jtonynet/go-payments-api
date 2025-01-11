import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._

class Payment extends Simulation {

  val httpProtocol = http
    .baseUrl("http://transaction-rest:8080")
    .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .acceptEncodingHeader("gzip, deflate")
    .userAgentHeader("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:16.0) Gecko/20100101 Firefox/16.0")

  val paymentExecute = 
    feed(tsv("transactions.tsv").circular())
    .exec(http("Requisição para /payment")
      .post("/payment")
      .body(StringBody("#{payload}"))
      .header("Content-Type", "application/json")
      .check(
        status.is(200),
        jsonPath("$.code").transform(code => if (code == "07") throw new Exception("General Error: code 07") else code)
      )
    )

  val testPaymentExecute = scenario("Test Payments").exec(paymentExecute)

  private val tps = 25
  private val window = 5.minutes
  private val activeUsers = (window.toSeconds * tps).toInt

  setUp(
    testPaymentExecute.inject(rampUsers(activeUsers).during(window))
  ).protocols(httpProtocol)
}
