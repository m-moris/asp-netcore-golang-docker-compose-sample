using System.Threading.Tasks;
using Grpc.Net.Client;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using WebApplication.HelloWorld;

namespace WebApplication.Controllers
{
    [ApiController]
    [Route("")]
    public class HelloController : ControllerBase
    {
        private readonly ILogger<HelloController> _logger;
        private readonly string _address;

        public HelloController(IOptions<GrpcOptions> opt, ILogger<HelloController> logger)
        {
            _logger = logger;
            _address = $"http://{opt.Value.Server}:{opt.Value.Port}";
            _logger.LogInformation($"{_address}");
        }

        [HttpGet("hello")]
        public async Task<IActionResult> Hello()
        {
            _logger.LogInformation("Hello");
            var channel = GrpcChannel.ForAddress(_address);
            var client = new Greeter.GreeterClient(channel);

            var req = new HelloRequest {Name = "Azure"};
            var res = await client.SayHelloAsync(req);

            return new JsonResult(
                new NetCoreHelloResponse
                {
                    Message = res.Message,
                    Date = res.Date,
                });
        }
    }

    public class NetCoreHelloResponse
    {
        public string Message { get; set; }
        public string Date { get; set; }
    }
}