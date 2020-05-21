package org.jmoney8080.search

import io.dropwizard.Application
import io.dropwizard.Configuration
import io.dropwizard.setup.Environment
import org.elasticsearch.action.search.SearchRequest
import org.elasticsearch.action.search.SearchResponse
import org.elasticsearch.search.builder.SearchSourceBuilder
import java.util.*
import javax.ws.rs.GET
import javax.ws.rs.Path
import javax.ws.rs.Produces
import javax.ws.rs.QueryParam
import javax.ws.rs.core.MediaType
import java.util.concurrent.TimeUnit
import org.elasticsearch.common.unit.TimeValue
import org.elasticsearch.index.query.QueryBuilders
import java.net.InetAddress
import org.elasticsearch.common.transport.TransportAddress
import org.elasticsearch.client.transport.TransportClient
import org.elasticsearch.common.settings.Settings
import org.elasticsearch.transport.client.PreBuiltTransportClient


class SearchServiceConfig(val name: String = "unknown") : Configuration()

fun main(args: Array<String>) {
    SearchServiceApp().run(*args) // use collection as a varargs
}

class SearchServiceApp : Application<SearchServiceConfig>() {
    override fun run(configuration: SearchServiceConfig, environment: Environment) {

        val client = PreBuiltTransportClient(Settings.builder().put("cluster.name", "docker-cluster").build())
            .addTransportAddress(TransportAddress(InetAddress.getByName("localhost"), 9300))

        val searchServiceComponent = SearchServiceComponent(client)

        environment.jersey().register(searchServiceComponent)
    }
}

@Path("/")
@Produces(MediaType.APPLICATION_JSON)
class SearchServiceComponent constructor(val esClient: TransportClient) {

    @Path("/search")
    @GET
    fun search(@QueryParam("agentid") agentid: String): SearchResponse {

        val sourceBuilder = SearchSourceBuilder()
        sourceBuilder.query(QueryBuilders.matchQuery("AgentID", agentid))
        sourceBuilder.timeout(TimeValue(60, TimeUnit.SECONDS))
        val reqeust = SearchRequest("datasets").source(sourceBuilder)
        return esClient.search(reqeust).actionGet()
    }
}

data class DataSearch(val agentId: String, val datasetId: String) {}
