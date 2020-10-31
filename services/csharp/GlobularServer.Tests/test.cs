using System;
using Xunit;
using Globular;

namespace GlobularServer.Tests
{
    public class UnitTest1
    {
        [Fact]
        public void TestCreateService()
        {
            // Test create service instance...
            GlobularService service = new GlobularService();
            // initialyse it...
            service.init();
            Assert.Equal(service.getPath(), "E:\\Project\\src\\github.com\\davecourtois\\Globular\\csharp\\GlobularServer.Tests\\bin\\Debug\netcoreapp2.1\\GlobularServer.Tests.dll");
        }
    }
}
