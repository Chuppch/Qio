package com.chuppch.config;

import com.zaxxer.hikari.HikariDataSource;
import org.apache.ibatis.session.SqlSessionFactory;
import org.mybatis.spring.SqlSessionFactoryBean;
import org.mybatis.spring.SqlSessionTemplate;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;
import org.springframework.core.io.support.PathMatchingResourcePatternResolver;
import org.springframework.jdbc.core.JdbcTemplate;

import javax.sql.DataSource;

/**
 * @author chuppch
 * @description
 * @create 2025/12/16
 */
@Configuration
public class DataSourceConfig {

    // ============================== Mysql 数据库配置 =========================================
    @Bean("mysqlDataSource")
    @Primary
    public DataSource mysqlDataSource(@Value("${spring.datasource.mysql.driver-class-name}") String driverClassName,
                                      @Value("${spring.datasource.mysql.url}") String url,
                                      @Value("${spring.datasource.mysql.username}") String username,
                                      @Value("${spring.datasource.mysql.password}") String password,
                                      @Value("${spring.datasource.mysql.hikari.maximum-pool-size:10}") int maximumPoolSize,
                                      @Value("${spring.datasource.mysql.hikari.minimum-idle:5}") int minimumIdle,
                                      @Value("${spring.datasource.mysql.hikari.idle-timeout:30000}") long idleTimeout,
                                      @Value("${spring.datasource.mysql.hikari.connection-timeout:30000}") long connectionTimeout,
                                      @Value("${spring.datasource.mysql.hikari.max-lifetime:1800000}") long maxLifetime,
                                      @Value("${spring.datasource.mysql.hikari.pool-name:Retail_HikariCP}") String poolName,
                                      @Value("${spring.datasource.mysql.hikari.auto-commit:true}") boolean autoCommit,
                                      @Value("${spring.datasource.mysql.hikari.connection-test-query:SELECT 1}") String connectionTestQuery) {
        // 连接池配置
        HikariDataSource dataSource = new HikariDataSource();
        dataSource.setDriverClassName(driverClassName);
        dataSource.setJdbcUrl(url);
        dataSource.setUsername(username);
        dataSource.setPassword(password);

        dataSource.setMaximumPoolSize(maximumPoolSize);
        dataSource.setMinimumIdle(minimumIdle);
        dataSource.setIdleTimeout(idleTimeout);
        dataSource.setConnectionTimeout(connectionTimeout);
        dataSource.setMaxLifetime(maxLifetime);
        dataSource.setPoolName(poolName);
        dataSource.setAutoCommit(autoCommit);
        dataSource.setConnectionTestQuery(connectionTestQuery);

        return dataSource;
    }

    // 配置Mybatis
    @Bean("sqlSessionFactory")
    public SqlSessionFactory sqlSessionFactory(@Qualifier("mysqlDataSource") DataSource mysqlDataSource) throws Exception {
        SqlSessionFactoryBean sqlSessionFactoryBean = new SqlSessionFactoryBean();
        sqlSessionFactoryBean.setDataSource(mysqlDataSource);

        //设置MyBatis配置文件位置
        PathMatchingResourcePatternResolver resolver = new PathMatchingResourcePatternResolver();
        sqlSessionFactoryBean.setConfigLocation(resolver.getResource("classpath:mybatis/config/mybatis-config.xml"));

        // 设置Mapper XML 文件位置
        sqlSessionFactoryBean.setMapperLocations(resolver.getResources("classpath:mybatis/mapper/*.xml"));
        return sqlSessionFactoryBean.getObject();
    }

    // 配置MySQL客户端
    @Bean("sqlSessionTemplate")
    public SqlSessionTemplate sqlSessionTemplate(@Qualifier("sqlSessionFactory") SqlSessionFactory sqlSessionFactory) {
        return new SqlSessionTemplate(sqlSessionFactory);
    }


    // ============================== PostgreSQL 数据库配置 =========================================

    @Bean("pgVectorDataSource")
    public DataSource pgVectorDataSource(@Value("${spring.datasource.pgvector.driver-class-name}") String driverClassName,
                                         @Value("${spring.datasource.pgvector.url}") String url,
                                         @Value("${spring.datasource.pgvector.username}") String username,
                                         @Value("${spring.datasource.pgvector.password}") String password,
                                         @Value("${spring.datasource.pgvector.hikari.maximum-pool-size:5}") int maximumPoolSize,
                                         @Value("${spring.datasource.pgvector.hikari.minimum-idle:2}") int minimumIdle,
                                         @Value("${spring.datasource.pgvector.hikari.idle-timeout:30000}") long idleTimeout,
                                         @Value("${spring.datasource.pgvector.hikari.connection-timeout:30000}") long connectionTimeout,
                                         @Value("${spring.datasource.pgvector.hikari.max-lifetime:1800000}") long maxLifetime,
                                         @Value("${spring.datasource.pgvector.hikari.pool-name:PgVectorHikariPool}") String poolName,
                                         @Value("${spring.datasource.pgvector.hikari.auto-commit:true}") boolean autoCommit,
                                         @Value("${spring.datasource.pgvector.hikari.connection-test-query:SELECT 1}") String connectionTestQuery,
                                         @Value("${spring.datasource.pgvector.hikari.validation-timeout:5000}") long validationTimeout,
                                         @Value("${spring.datasource.pgvector.hikari.leak-detection-threshold:0}") long leakDetectionThreshold) {

        // 连接池配置
        HikariDataSource dataSource = new HikariDataSource();
        dataSource.setDriverClassName(driverClassName);
        dataSource.setJdbcUrl(url);
        dataSource.setUsername(username);
        dataSource.setPassword(password);

        dataSource.setMaximumPoolSize(maximumPoolSize);
        dataSource.setMinimumIdle(minimumIdle);
        dataSource.setIdleTimeout(idleTimeout);
        dataSource.setConnectionTimeout(connectionTimeout);
        dataSource.setMaxLifetime(maxLifetime);
        dataSource.setPoolName(poolName);
        dataSource.setAutoCommit(autoCommit);
        dataSource.setConnectionTestQuery(connectionTestQuery);
        // 连接验证超时：验证连接是否有效时的超时时间
        dataSource.setValidationTimeout(validationTimeout);
        // 连接泄漏检测：如果连接在指定时间内未归还，记录警告日志（0表示禁用）
        if (leakDetectionThreshold > 0) {
            dataSource.setLeakDetectionThreshold(leakDetectionThreshold);
        }
        // initializationFailTimeout: -1 表示不超时，0 表示立即失败，>0 表示超时时间（毫秒）
        // 设置为 1 秒（1000ms），如果连接失败则快速失败，避免启动时长时间等待
        dataSource.setInitializationFailTimeout(1000);

        return dataSource;
    }

    @Bean("pgVectorJdbcTemplate")
    public JdbcTemplate pgVectorJdbcTemplate(@Qualifier("pgVectorDataSource") DataSource pgVectorDataSource) {
        return new JdbcTemplate(pgVectorDataSource);
    }
}